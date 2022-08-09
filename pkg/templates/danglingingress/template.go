package danglingingress

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/danglingingress/internal/params"
	v1 "k8s.io/api/core/v1"
	networkingV1 "k8s.io/api/networking/v1"
)

const (
	templateKey = "dangling-ingress"
)

func getSelectorsFromIngressBackend(b *networkingV1.IngressBackend) string {
	service := b.Service
	if service == nil {
		return ""
	}

	return service.Name
}

func getSelectorsFromIngress(ingress *networkingV1.Ingress) map[string]bool {
	selectors := map[string]bool{}

	if defaultBack := ingress.Spec.DefaultBackend; defaultBack != nil {
		if s := getSelectorsFromIngressBackend(defaultBack); s != "" {
			selectors[s] = true
		}
	}

	for _, r := range ingress.Spec.Rules {
		spec := r.IngressRuleValue.HTTP
		if spec == nil {
			continue
		}

		for _, p := range spec.Paths {
			if s := getSelectorsFromIngressBackend(&p.Backend); s != "" {
				selectors[s] = true
			}
		}
	}

	return selectors
}

func init() {
	templates.Register(check.Template{
		HumanName:   "Dangling Ingress",
		Key:         templateKey,
		Description: "Flag ingress which do not match any service",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.Service},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(_ params.Params) (check.Func, error) {
			return func(lintCtx lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				ingress, ok := object.K8sObject.(*networkingV1.Ingress)
				if !ok {
					return nil
				}

				selectors := getSelectorsFromIngress(ingress)

				// if there aren't any service selectors found we assume that the backend's are specified
				// by resources and will skip over this.
				if len(selectors) == 0 {
					return nil
				}

				for _, obj := range lintCtx.Objects() {
					if ingress.Namespace != obj.K8sObject.GetNamespace() {
						continue
					}

					service, ok := obj.K8sObject.(*v1.Service)
					if !ok {
						continue
					}

					serviceName := service.ObjectMeta.Name
					if _, ok := selectors[serviceName]; ok {
						delete(selectors, serviceName)
						if len(selectors) == 0 {
							// Found the all!
							return nil
						}
					}
				}

				var dig []diagnostic.Diagnostic

				for k := range selectors {
					dig = append(dig, diagnostic.Diagnostic{
						Message: fmt.Sprintf("no service found matching ingress labels (%s)", k),
					})
				}

				return dig
			}, nil
		}),
	})
}
