package readsecret

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/readsecret/internal/params"
	"golang.stackrox.io/kube-linter/pkg/templates/util"
	v1 "k8s.io/api/core/v1"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Read Secret From Environment Variables",
		Key:         "read-secret-from-env-var",
		Description: "Flag environment variables that use SecretKeyRef",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(_ params.Params) (check.Func, error) {
			return util.PerContainerCheck(func(container *v1.Container) []diagnostic.Diagnostic {
				var results []diagnostic.Diagnostic
				for _, envVar := range container.Env {
					if envVar.ValueFrom != nil && envVar.ValueFrom.SecretKeyRef != nil {
						results = append(results, diagnostic.Diagnostic{
							Message: fmt.Sprintf("environment variable %q in container %q uses SecretKeyRef", envVar.Name, container.Name),
						})
					}
				}
				return results
			}), nil
		}),
	})
}
