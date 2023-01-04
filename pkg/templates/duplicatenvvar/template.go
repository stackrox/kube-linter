package duplicateenvvar

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/duplicatenvvar/internal/params"
	"golang.stackrox.io/kube-linter/pkg/templates/util"
	v1 "k8s.io/api/core/v1"
)

const templateKey = "duplicate-env-var"

func init() {
	templates.Register(check.Template{
		HumanName:   "Duplicate Environment Variables",
		Key:         templateKey,
		Description: "Flag Duplicate Env Variables names",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(_ params.Params) (check.Func, error) {
			return util.PerContainerCheck(func(container *v1.Container) []diagnostic.Diagnostic {
				var results []diagnostic.Diagnostic
				envVarNames := map[string]int{}

				for _, envVar := range container.Env {
					// Ensure we only report on error per env var
					envVarNames[envVar.Name]++
					if num, ok := envVarNames[envVar.Name]; !ok || num != 2 {
						continue
					}
					results = append(results, diagnostic.Diagnostic{
						Message: fmt.Sprintf(
							"Duplicate environment variable %s in container %q found",
							envVar.Name,
							container.Name,
						),
					})
				}
				return results
			}), nil
		}),
	})
}
