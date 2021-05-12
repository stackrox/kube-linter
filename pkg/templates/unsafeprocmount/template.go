package unsafeprocmount

import (
	"fmt"
	"strings"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/unsafeprocmount/internal/params"
	"golang.stackrox.io/kube-linter/pkg/templates/util"
	v1 "k8s.io/api/core/v1"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Unsafe Proc Mount",
		Key:         "unsafe-proc-mount",
		Description: "Flag containers of unsafe proc mount",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(_ params.Params) (check.Func, error) {
			return util.PerContainerCheck(func(container *v1.Container) []diagnostic.Diagnostic {
				if container.SecurityContext != nil && container.SecurityContext.ProcMount != nil {
					if strings.EqualFold(string(*container.SecurityContext.ProcMount), "Unmasked") {
						return []diagnostic.Diagnostic{{Message: fmt.Sprintf("container %q exposes /proc unsafely (via procMount=Unmasked).", container.Name)}}
					}
				}
				return nil
			}), nil
		}),
	})
}
