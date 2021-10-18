package danglingnetworkpolicy

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/danglingnetworkpolicy/internal/params"

	v1 "k8s.io/api/networking/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	templateKey = "dangling-networkpolicy"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Dangling NetworkPolicies",
		Key:         "dangling-networkpolicy",
		Description: "Flag NetworkPolicies which do not match any application",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(_ params.Params) (check.Func, error) {
			return func(lintCtx lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				networkpolicy, ok := object.K8sObject.(*v1.NetworkPolicy)
				if !ok {
					return nil
				}
				podselector := networkpolicy.Spec.PodSelector
				labelSelector, err := metaV1.LabelSelectorAsSelector(&metaV1.LabelSelector{MatchLabels: podselector.MatchLabels,
					MatchExpressions: podselector.MatchExpressions})
				if err != nil {
					return []diagnostic.Diagnostic{{
						Message: fmt.Sprintf("networkpolicy has invalid podSelector: %v", err),
					}}
				}
				for _, obj := range lintCtx.Objects() {
					podTemplateSpec, hasPods := extract.PodTemplateSpec(obj.K8sObject)
					if !hasPods {
						continue
					}
					if networkpolicy.Namespace != obj.K8sObject.GetNamespace() {
						continue
					}
					if labelSelector.Matches(labels.Set(podTemplateSpec.Labels)) {
						// Found!
						return nil
					}
				}
				return []diagnostic.Diagnostic{{Message: fmt.Sprintf("no pods found matching networkpolicy's podSelector labels (%v) ", podselector)}}
			}, nil
		}),
	})
}
