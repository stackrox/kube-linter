package servicetype

import (
	"fmt"
	"strings"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/servicetype/internal/params"
	v1 "k8s.io/api/core/v1"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Forbidden Service Types",
		Key:         "forbidden-service-types",
		Description: "Flag forbidden services",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.Service},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			return func(lintCtx lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {

				service, ok := object.K8sObject.(*v1.Service)
				if !ok {
					return nil
				}
				var results []diagnostic.Diagnostic
				for _, servicetype := range p.ForbiddenServiceTypes {
					if strings.EqualFold(string(service.Spec.Type), servicetype) {
						results = append(results, diagnostic.Diagnostic{Message: fmt.Sprintf("%q service type is forbidden.", servicetype)})
					}
				}
				return results
			}, nil
		}),
	})
}
