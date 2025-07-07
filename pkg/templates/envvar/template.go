package envvar

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/matcher"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/envvar/internal/params"
	"golang.stackrox.io/kube-linter/pkg/templates/util"
	v1 "k8s.io/api/core/v1"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Environment Variables",
		Key:         "env-var",
		Description: "Flag environment variables that match the provided patterns",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			nameMatcher, err := matcher.ForString(p.Name)
			if err != nil {
				return nil, fmt.Errorf("invalid name: %w", err)
			}
			valueMatcher, err := matcher.ForString(p.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid value: %w", err)
			}
			return util.PerContainerCheck(func(container *v1.Container) []diagnostic.Diagnostic {
				var results []diagnostic.Diagnostic
				for _, envVar := range container.Env {
					if nameMatcher(envVar.Name) && valueMatcher(envVar.Value) {
						results = append(results, diagnostic.Diagnostic{
							Message: fmt.Sprintf("environment variable %s in container %q found", envVar.Name, container.Name),
						})
					}
				}
				return results
			}), nil
		}),
	})
}
