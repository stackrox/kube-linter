package privilegeescalation

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/privilegeescalation/internal/params"
	"golang.stackrox.io/kube-linter/pkg/templates/util"
	v1 "k8s.io/api/core/v1"
)

const (
	sysAdminCapability = "SYS_ADMIN"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Privilege Escalation on Containers",
		Key:         "privilege-escalation-container",
		Description: "Flag containers of allowing privilege escalation",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(_ params.Params) (check.Func, error) {
			return util.PerContainerCheck(func(container *v1.Container) []diagnostic.Diagnostic {
				securityContext := container.SecurityContext
				if securityContext == nil {
					return nil
				}
				if securityContext.AllowPrivilegeEscalation != nil && *securityContext.AllowPrivilegeEscalation {
					return []diagnostic.Diagnostic{{Message: fmt.Sprintf("container %q has AllowPrivilegeEscalation set to true.", container.Name)}}
				}
				if securityContext.Privileged != nil && *securityContext.Privileged {
					return []diagnostic.Diagnostic{{Message: fmt.Sprintf("container %q is Privileged hence allows privilege escalation.", container.Name)}}
				}
				if securityContext.Capabilities != nil {
					for _, capability := range securityContext.Capabilities.Add {
						if capability == sysAdminCapability {
							return []diagnostic.Diagnostic{{Message: fmt.Sprintf("container %q has SYS_ADMIN capability hence allows privilege escalation.", container.Name)}}
						}
					}
				}
				return nil
			}), nil
		}),
	})
}
