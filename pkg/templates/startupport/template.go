package startupport

import (
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/startupport/internal/params"
	"golang.stackrox.io/kube-linter/pkg/templates/util"
	v1 "k8s.io/api/core/v1"
)

const templateKey = "startup-port"

func init() {
	templates.Register(check.Template{
		HumanName:   "Startup Port Exposed",
		Key:         templateKey,
		Description: "Flag containers with an Startup probe to not exposed port.",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(_ params.Params) (check.Func, error) {
			return util.PerNonInitContainerCheck(func(container *v1.Container) []diagnostic.Diagnostic {
				return util.CheckProbePort(container, container.StartupProbe)
			}), nil
		}),
	})
}
