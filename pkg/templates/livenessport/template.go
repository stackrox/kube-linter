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

var sentinel = struct{}{}

func init() {
	templates.Register(check.Template{
		HumanName:   "Liveness Port Not Open",
		Key:         templateKey,
		Description: "Flag containers with an HTTP liveness probe to not exposed port.",
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

				ports := map[intstr.IntOrString]struct{}{}
				for _, port := range container.Ports {
					if port.Protocol != "" && port.Protocol != v1.ProtocolTCP {
						continue
					}
					ports[intstr.FromInt32(port.ContainerPort)] = sentinel
					ports[intstr.FromString(port.Name)] = sentinel
				}

				if httpProbe := container.LivenessProbe.ProbeHandler.HTTPGet; httpProbe != nil {
					if _, ok := ports[httpProbe.Port]; !ok {
						return []diagnostic.Diagnostic{{
							Message: fmt.Sprintf("container %q does not expose port %s for the HTTPGet", container.Name, httpProbe.Port.String()),
						}}
					}
				}

				if tcpProbe := container.LivenessProbe.ProbeHandler.TCPSocket; tcpProbe != nil {
					if _, ok := ports[tcpProbe.Port]; !ok {
						return []diagnostic.Diagnostic{{
							Message: fmt.Sprintf("container %q does not expose port %s for the TCPSocket", container.Name, tcpProbe.Port.String()),
						}}
					}
				}
				return nil
			}), nil
		}),
	})
}
