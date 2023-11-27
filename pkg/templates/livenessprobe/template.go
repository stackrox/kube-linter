package livenessprobe

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/livenessprobe/internal/params"
	"golang.stackrox.io/kube-linter/pkg/templates/util"
	v1 "k8s.io/api/core/v1"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Liveness Probe Not Specified",
		Key:         "liveness-probe",
		Description: "Flag containers that don't specify a liveness probe",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters: params.ParamDescs,

		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(_ params.Params) (check.Func, error) {
			return util.PerNonInitContainerCheck(func(container *v1.Container) []diagnostic.Diagnostic {
				if container.LivenessProbe == nil {
					return []diagnostic.Diagnostic{{Message: fmt.Sprintf("container %q does not specify a liveness probe", container.Name)}}
				}
				return nil
			}), nil
		}),
	})
}
