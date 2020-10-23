package privileged

import (
	"fmt"

	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/diagnostic"
	"golang.stackrox.io/kube-linter/internal/objectkinds"
	"golang.stackrox.io/kube-linter/internal/templates"
	"golang.stackrox.io/kube-linter/internal/templates/privileged/internal/params"
	"golang.stackrox.io/kube-linter/internal/templates/util"
	v1 "k8s.io/api/core/v1"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Privileged Containers",
		Key:         "privileged",
		Description: "Flag privileged containers",
		SupportedObjectKinds: check.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(_ params.Params) (check.Func, error) {
			return util.PerContainerCheck(func(container *v1.Container) []diagnostic.Diagnostic {
				if securityContext := container.SecurityContext; securityContext != nil {
					if securityContext.Privileged != nil && *securityContext.Privileged {
						return []diagnostic.Diagnostic{{Message: fmt.Sprintf("container %q is privileged", container.Name)}}
					}
				}
				return nil
			}), nil
		}),
	})
}
