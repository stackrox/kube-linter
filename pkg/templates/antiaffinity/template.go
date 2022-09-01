package antiaffinity

import (
	"fmt"
	"strings"

	"golang.stackrox.io/kube-linter/internal/stringutils"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/matcher"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/antiaffinity/internal/params"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	templateKey = "anti-affinity"
)

func defaultTopologyKeyMatcher(topologyKey string) bool {
	return topologyKey == "kubernetes.io/hostname"
}

func init() {
	templates.Register(check.Template{
		HumanName: "Anti affinity not specified",
		Key:       templateKey,
		Description: "Flag objects with multiple replicas but inter-pod anti affinity not specified in the pod " +
			"template spec",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			var topologyKeyMatcher func(string) bool
			if p.TopologyKey == "" {
				topologyKeyMatcher = defaultTopologyKeyMatcher
			} else {
				var err error
				topologyKeyMatcher, err = matcher.ForString(p.TopologyKey)
				if err != nil {
					return nil, err
				}
			}
			return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				replicas, found := extract.Replicas(object.K8sObject)
				if !found {
					return nil
				}
				if int(replicas) < p.MinReplicas {
					return nil
				}
				podTemplateSpec, hasPods := extract.PodTemplateSpec(object.K8sObject)
				if !hasPods {
					return nil
				}
				namespace := object.K8sObject.GetNamespace()
				affinity := podTemplateSpec.Spec.Affinity
				// Short-circuit if no affinity rule is specified within the pod spec.
				if affinity == nil || affinity.PodAntiAffinity == nil {
					return []diagnostic.Diagnostic{
						{Message: fmt.Sprintf("object has %d %s but does not specify inter pod anti-affinity",
							replicas, stringutils.Ternary(replicas > 1, "replicas", "replica"))},
					}
				}
				var foundIssues []diagnostic.Diagnostic
				preferredAffinity := affinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution
				requiredAffinity := affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution
				// Short-circuit if affinity rule is specified but both preferred and required are empty.
				if len(preferredAffinity) == 0 && len(requiredAffinity) == 0 {
					return []diagnostic.Diagnostic{
						{Message: fmt.Sprintf("object has %d %s but does not specify preferred or required "+
							"inter pod anti-affinity during scheduling",
							replicas, stringutils.Ternary(replicas > 1, "replicas", "replica"))},
					}
				}

				for _, preferred := range preferredAffinity {
					err := validateAffinityTermMatchesAgainstNodes(preferred.PodAffinityTerm,
						namespace, podTemplateSpec.Labels, topologyKeyMatcher)
					if err == nil {
						return nil
					}
					foundIssues = append(foundIssues, diagnostic.Diagnostic{
						Message: err.Error(),
					})
				}
				for _, required := range requiredAffinity {
					err := validateAffinityTermMatchesAgainstNodes(required, namespace,
						podTemplateSpec.Labels, topologyKeyMatcher)
					if err == nil {
						return nil
					}
					foundIssues = append(foundIssues, diagnostic.Diagnostic{
						Message: err.Error(),
					})
				}
				return foundIssues
			}, nil
		}),
	})
}

func validateAffinityTermMatchesAgainstNodes(affinityTerm coreV1.PodAffinityTerm, podNamespace string,
	podLabels map[string]string, topologyKeyMatcher func(string) bool) error {
	// If namespaces is not specified in the affinity term, that means the affinity term implicitly applies to
	// the pod's namespace.
	if len(affinityTerm.Namespaces) > 0 {
		var matchingNSFound bool
		for _, ns := range affinityTerm.Namespaces {
			if ns == podNamespace {
				matchingNSFound = true
				break
			}
		}
		if !matchingNSFound {
			return fmt.Errorf("pod's namespace %q not found in anti-affinity's namespaces [%s]",
				podNamespace, strings.Join(affinityTerm.Namespaces, ", "))
		}
	}
	labelSelector, err := metaV1.LabelSelectorAsSelector(affinityTerm.LabelSelector)
	if err != nil {
		return err
	}

	if !topologyKeyMatcher(affinityTerm.TopologyKey) {
		return fmt.Errorf("anti-affinity's topology key does not match %q", affinityTerm.TopologyKey)
	}
	if !labelSelector.Matches(labels.Set(podLabels)) {
		return fmt.Errorf("pod's labels %q do not match with anti-affinity's labels %q",
			labels.Set(podLabels).String(), labelSelector.String())
	}
	return nil
}
