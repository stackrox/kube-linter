package envvar

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/diagnostic"
	"golang.stackrox.io/kube-linter/internal/matcher"
	"golang.stackrox.io/kube-linter/internal/objectkinds"
	"golang.stackrox.io/kube-linter/internal/templates"
	"golang.stackrox.io/kube-linter/internal/templates/envvar/internal/params"
	"golang.stackrox.io/kube-linter/internal/templates/util"
	v1 "k8s.io/api/core/v1"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Environment Variables",
		Key:         "env-var",
		Description: "Flag environment variables that match the provided patterns",
		SupportedObjectKinds: check.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			nameMatcher, err := matcher.ForString(p.Name)
			if err != nil {
				return nil, errors.Wrap(err, "invalid name")
			}
			valueMatcher, err := matcher.ForString(p.Value)
			if err != nil {
				return nil, errors.Wrap(err, "invalid value")
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
