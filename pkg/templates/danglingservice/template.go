package danglingservice

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/danglingservice/internal/params"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const templateKey = "dangling-service"

func init() {
	templates.Register(check.Template{
		HumanName:   "Dangling Services",
		Key:         templateKey,
		Description: "Flag services which do not match any application",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			return func(lintCtx lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
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

				for _, ignoredLabel := range p.IgnoredLabels {
					delete(selector, ignoredLabel)
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
					if service.Namespace != obj.K8sObject.GetNamespace() {
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
