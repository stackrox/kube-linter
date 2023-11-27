package livenessport

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/livenessport/internal/params"
	"golang.stackrox.io/kube-linter/pkg/templates/util"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const templateKey = "liveness-http-port"

var sentinal = struct{}{}

func init() {
	templates.Register(check.Template{
		HumanName:   "Liveness Port Not Open",
		Key:         templateKey,
		Description: "Flag containers have a http liveness prop with a port they didn't open.",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(_ params.Params) (check.Func, error) {
			return util.PerNonInitContainerCheck(func(container *v1.Container) []diagnostic.Diagnostic {
				if container.LivenessProbe == nil {
					return nil
				}

				httpProbe := container.LivenessProbe.ProbeHandler.HTTPGet
				if httpProbe == nil {
					return nil
				}

				ports := map[intstr.IntOrString]struct{}{}
				for _, port := range container.Ports {
					ports[intstr.FromInt(int(port.ContainerPort))] = sentinal
					ports[intstr.FromString(port.Name)] = sentinal
				}

				if _, ok := ports[httpProbe.Port]; !ok {
					return []diagnostic.Diagnostic{{
						Message: fmt.Sprintf("container %q does not have an open port %s", container.Name, httpProbe.Port.String()),
					}}
				}
				return nil
			}), nil
		}),
	})
}
