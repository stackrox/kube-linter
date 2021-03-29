package antiaffinity

import (
	"fmt"

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
		HumanName:   "Anti affinity not specified",
		Key:         templateKey,
		Description: "Flag objects with multiple replicas but inter-pod anti affinity not specified in the pod template spec",
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
				if affinity := podTemplateSpec.Spec.Affinity; affinity != nil && affinity.PodAntiAffinity != nil {
					preferredAffinity := affinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution
					requiredAffinity := affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution
					for _, preferred := range preferredAffinity {
						if affinityTermMatchesLabelsAgainstNodes(preferred.PodAffinityTerm, podTemplateSpec.Namespace, podTemplateSpec.Labels, topologyKeyMatcher) {
							return nil
						}
					}
					for _, required := range requiredAffinity {
						if affinityTermMatchesLabelsAgainstNodes(required, podTemplateSpec.Namespace, podTemplateSpec.Labels, topologyKeyMatcher) {
							return nil
						}
					}
				}
				return []diagnostic.Diagnostic{
					{Message: fmt.Sprintf("object has %d %s but does not specify inter pod anti-affinity", replicas, stringutils.Ternary(replicas > 1, "replicas", "replica"))},
				}
			}, nil
		}),
	})
}

func affinityTermMatchesLabelsAgainstNodes(affinityTerm coreV1.PodAffinityTerm, podNamespace string, podLabels map[string]string, topologyKeyMatcher func(string) bool) bool {
	// If namespaces is not specified in the affinity term, that means the affinity term implicitly applies to the pod's namespace.
	if len(affinityTerm.Namespaces) > 0 {
		var matchingNSFound bool
		for _, ns := range affinityTerm.Namespaces {
			if ns == podNamespace {
				matchingNSFound = true
				break
			}
		}
		if !matchingNSFound {
			return false
		}
	}
	labelSelector, err := metaV1.LabelSelectorAsSelector(affinityTerm.LabelSelector)
	if err != nil {
		return false
	}
	if topologyKeyMatcher(affinityTerm.TopologyKey) && labelSelector.Matches(labels.Set(podLabels)) {
		return true
	}
	return false
}
