package danglingnetworkpolicypeer

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/danglingnetworkpolicypeer/internal/params"

	v1 "k8s.io/api/networking/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	templateKey = "dangling-networkpolicypeer-podselector"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Dangling NetworkPolicyPeer PodSelector",
		Key:         templateKey,
		Description: "Flag NetworkPolicyPeer in Ingress/Egress rules which their podselector do not match any application. Applied to peers consisting with podSelectors only.",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(_ params.Params) (check.Func, error) {
			return func(lintCtx lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				var results []diagnostic.Diagnostic
				networkpolicy, ok := object.K8sObject.(*v1.NetworkPolicy)
				if !ok {
					return nil
				}
				ingressRules := networkpolicy.Spec.Ingress
				for _, inrule := range ingressRules {
					for _, peer := range inrule.From {
						res := getAndCheckifPodSelectorMatchesPods(peer, lintCtx, networkpolicy.Namespace)
						if res != nil {
							results = append(results, res[len(res)-1])
						}
					}
				}
				egressRules := networkpolicy.Spec.Egress
				for _, torule := range egressRules {
					for _, peer := range torule.To {
						res := getAndCheckifPodSelectorMatchesPods(peer, lintCtx, networkpolicy.Namespace)
						if res != nil {
							results = append(results, res[len(res)-1])
						}
					}
				}
				return results
			}, nil
		}),
	})
}

func getAndCheckifPodSelectorMatchesPods(peer v1.NetworkPolicyPeer, lintCtx lintcontext.LintContext, currNamespace string) []diagnostic.Diagnostic {
	podSelector := peer.PodSelector
	if podSelector == nil {
		return nil
	}
	nsSelector := peer.NamespaceSelector
	if nsSelector != nil {
		return nil // For now, we assume all pods with namespace selectors are okay
	}
	return findMatchingPods(podSelector, lintCtx, currNamespace)
}

func findMatchingPods(podSelector *metaV1.LabelSelector, lintCtx lintcontext.LintContext, currNamespace string) []diagnostic.Diagnostic {
	labelSelector, err := metaV1.LabelSelectorAsSelector(&metaV1.LabelSelector{MatchLabels: podSelector.MatchLabels,
		MatchExpressions: podSelector.MatchExpressions})
	if err != nil {
		return []diagnostic.Diagnostic{{Message: fmt.Sprintf("networkpolicy ingress rule has invalid podSelector: %v", err)}}
	}
	for _, obj := range lintCtx.Objects() {
		podTemplateSpec, hasPods := extract.PodTemplateSpec(obj.K8sObject)
		if !hasPods {
			continue
		}
		if currNamespace != obj.K8sObject.GetNamespace() {
			continue
		}
		if labelSelector.Matches(labels.Set(podTemplateSpec.Labels)) {
			return nil // found
		}
	}
	return []diagnostic.Diagnostic{{Message: fmt.Sprintf("no pods found matching networkpolicy rule's podSelector labels (%v)", labelSelector)}}
}
