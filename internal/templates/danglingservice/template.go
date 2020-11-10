package danglingservice

import (
	"fmt"

	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/diagnostic"
	"golang.stackrox.io/kube-linter/internal/extract"
	"golang.stackrox.io/kube-linter/internal/lintcontext"
	"golang.stackrox.io/kube-linter/internal/objectkinds"
	"golang.stackrox.io/kube-linter/internal/templates"
	"golang.stackrox.io/kube-linter/internal/templates/danglingservice/internal/params"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Dangling Services",
		Key:         "dangling-service",
		Description: "Flag services which do not match any application",
		SupportedObjectKinds: check.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(_ params.Params) (check.Func, error) {
			return func(lintCtx *lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				service, ok := object.K8sObject.(*v1.Service)
				if !ok {
					return nil
				}
				// Selector doesn't apply to external names.
				if service.Spec.Type == v1.ServiceTypeExternalName {
					return nil
				}
				selector := service.Spec.Selector
				if len(selector) == 0 {
					return []diagnostic.Diagnostic{{
						Message: "service has no selector specified",
					}}
				}
				labelSelector, err := metaV1.LabelSelectorAsSelector(&metaV1.LabelSelector{MatchLabels: selector})
				if err != nil {
					return []diagnostic.Diagnostic{{
						Message: fmt.Sprintf("service has invalid label selector: %v", err),
					}}
				}
				for _, obj := range lintCtx.Objects() {
					podTemplateSpec, hasPods := extract.PodTemplateSpec(obj.K8sObject)
					if !hasPods {
						continue
					}
					if service.Namespace != podTemplateSpec.Namespace {
						continue
					}
					if labelSelector.Matches(labels.Set(podTemplateSpec.Labels)) {
						// Found!
						return nil
					}
				}
				return []diagnostic.Diagnostic{{Message: fmt.Sprintf("no pods found matching service labels (%v)", selector)}}
			}, nil
		}),
	})
}
