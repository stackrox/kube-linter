package accesstoresources

import (
	"fmt"
	"regexp"
	"strings"

	"golang.stackrox.io/kube-linter/internal/stringutils"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/accesstoresources/internal/params"
	rbacV1 "k8s.io/api/rbac/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	templateKey = "access-to-resources"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Access to Resources",
		Key:         templateKey,
		Description: "Flag cluster role bindings and role bindings that grant access to the specified resource kinds and verbs",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{
				objectkinds.Role,
				objectkinds.ClusterRole,
				objectkinds.ClusterRoleBinding,
				objectkinds.RoleBinding},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			resourceRegexes := make([]*regexp.Regexp, 0, len(p.Resources))
			for _, res := range p.Resources {
				r, err := regexp.Compile(res)
				if err != nil {
					return nil, fmt.Errorf("invalid regex %s: %w", res, err)
				}
				resourceRegexes = append(resourceRegexes, r)
			}
			verbRegexes := make([]*regexp.Regexp, 0, len(p.Verbs))
			for _, verb := range p.Verbs {
				v, err := regexp.Compile(verb)
				if err != nil {
					return nil, fmt.Errorf("invalid regex %s: %w", verb, err)
				}
				verbRegexes = append(verbRegexes, v)
			}
			return func(lintCtx lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				rbinding, ok := object.K8sObject.(*rbacV1.RoleBinding)
				if ok {
					namespace := stringutils.OrDefault(rbinding.Namespace, "default")
					return findRole(rbinding.RoleRef.Name, namespace, lintCtx, resourceRegexes, verbRegexes, p.FlagRolesNotFound)
				}
				crbinding, ok := object.K8sObject.(*rbacV1.ClusterRoleBinding)
				if ok {
					return findClusterRole(crbinding.RoleRef.Name, lintCtx, resourceRegexes, verbRegexes, p.FlagRolesNotFound)
				}
				return nil
			}, nil
		}),
	})
}

// find clusterrole by name, and check if it has access to the specified resource kinds and verbs
func findClusterRole(name string, lintCtx lintcontext.LintContext, resourceRegexes, verbRegexes []*regexp.Regexp, flag bool) []diagnostic.Diagnostic {
	var results []diagnostic.Diagnostic
	var clusterroles []*rbacV1.ClusterRole
	for _, object := range lintCtx.Objects() {
		r, ok := object.K8sObject.(*rbacV1.ClusterRole)
		if !ok {
			continue
		}
		clusterroles = append(clusterroles, r)
	}

	var roleExists bool
	for _, r := range clusterroles {
		if r.Name == name && !strings.EqualFold(r.Name, "cluster_admin") {
			roleExists = true
			accesses := checkAccess(r.Rules, resourceRegexes, verbRegexes)
			if len(accesses) > 0 {
				results = append(results, diagnostic.Diagnostic{Message: fmt.Sprintf("binding to %q clusterrole that has %s", r.Name, strings.Join(accesses, ", "))})
			}
			if r.AggregationRule != nil && len(r.AggregationRule.ClusterRoleSelectors) > 0 {
				resultsAggregated := findAggregatedAccesses(clusterroles, r.AggregationRule.ClusterRoleSelectors, resourceRegexes, verbRegexes)
				results = append(results, resultsAggregated...)
			}
		}
	}
	if !roleExists && flag {
		results = append(results, diagnostic.Diagnostic{Message: fmt.Sprintf("clusterrole %q not found", name)})
	}
	return results
}

// find clusterroles by label selectors, and check if they have access to the specified resources and verbs
func findAggregatedAccesses(clusterroles []*rbacV1.ClusterRole, selectors []metaV1.LabelSelector, resourceRegexes, verbRegexes []*regexp.Regexp) []diagnostic.Diagnostic {
	var results []diagnostic.Diagnostic
	for _, s := range selectors {
		labelSelector, err := metaV1.LabelSelectorAsSelector(&metaV1.LabelSelector{MatchLabels: s.MatchLabels})
		if err != nil {
			continue
		}
		for _, r := range clusterroles {
			if labelSelector.Matches(labels.Set(r.GetLabels())) { // Found the aggregated clusterrole!
				accesses := checkAccess(r.Rules, resourceRegexes, verbRegexes)
				if len(accesses) > 0 {
					results = append(results, diagnostic.Diagnostic{Message: fmt.Sprintf("binding via aggregationRule to %q clusterrole that has %s", r.Name, strings.Join(accesses, ", "))})
				}
			}
		}
	}
	return results
}

// find role by name and namespace that has access to the specified resources and verbs
func findRole(name, namespace string, lintCtx lintcontext.LintContext, resources, verbs []*regexp.Regexp, flag bool) []diagnostic.Diagnostic {
	results := []diagnostic.Diagnostic{}
	var roleExists bool
	for _, object := range lintCtx.Objects() {
		r, ok := object.K8sObject.(*rbacV1.Role)
		if !ok {
			continue
		}
		ns := stringutils.OrDefault(r.Namespace, "default")
		if r.Name == name && ns == namespace {
			roleExists = true
			accesses := checkAccess(r.Rules, resources, verbs)
			if len(accesses) > 0 {
				results = append(results, diagnostic.Diagnostic{Message: fmt.Sprintf("binding to %q role that has %s", r.Name, strings.Join(accesses, ", "))})
			}
		}
	}
	if !roleExists && flag {
		results = append(results, diagnostic.Diagnostic{Message: fmt.Sprintf("role %q in namespace %q not found", name, namespace)})
	}
	return results
}

// find access verbs to a given resource kind
func checkAccess(rules []rbacV1.PolicyRule, resourceRegex, verbRegex []*regexp.Regexp) []string {
	var accesses []string
	for _, rule := range rules {
		var resources []string
		for _, res := range rule.Resources {
			if isInList(resourceRegex, res) {
				resources = append(resources, res)
			}
		}
		if len(resources) == 0 {
			continue
		}
		var verbs []string
		for _, verb := range rule.Verbs {
			if isInList(verbRegex, verb) {
				verbs = append(verbs, verb)
			}
		}
		if len(verbs) == 0 {
			continue
		}
		accesses = append(accesses, fmt.Sprintf("%v access to %v", verbs, resources))
	}
	return accesses
}

// isInList returns true if a match found in the list for the given name or a wildcard
func isInList(regexlist []*regexp.Regexp, name string) bool {
	for _, regex := range regexlist {
		if name == "*" || regex.MatchString(name) {
			return true
		}
	}
	return false
}
