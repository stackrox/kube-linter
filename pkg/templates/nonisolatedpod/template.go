package nonisolatedpod

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/nonisolatedpod/internal/params"
	networkingV1 "k8s.io/api/networking/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	templateKey = "non-isolated-pod"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Non Isolated Pods",
		Key:         templateKey,
		Description: "Flag Pod that is not selected by any networkPolicy",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.NetworkPolicy},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(_ params.Params) (check.Func, error) {
			return func(lintCtx lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				podTemplateSpec, found := extract.PodTemplateSpec(object.K8sObject)
				if !found {
					return nil
				}
				for _, obj := range lintCtx.Objects() {
					networkpolicy, ok := obj.K8sObject.(*networkingV1.NetworkPolicy)
					if !ok {
						continue
					}
					if object.K8sObject.GetNamespace() != obj.K8sObject.GetNamespace() {
						continue
					}
					podselector := networkpolicy.Spec.PodSelector
					labelSelector, err := metaV1.LabelSelectorAsSelector(&metaV1.LabelSelector{MatchLabels: podselector.MatchLabels,
						MatchExpressions: podselector.MatchExpressions})
					if err != nil {
						return []diagnostic.Diagnostic{{
							Message: fmt.Sprintf("networkpolicy has invalid podSelector: %v", err),
						}}
					}
					if labelSelector.Matches(labels.Set(podTemplateSpec.Labels)) {
						// Found!
						return nil
					}
				}
				return []diagnostic.Diagnostic{{
					Message: "pods created by this object are non-isolated",
				}}
			}, nil
		}),
	})
}
