package probeport

import (
	"fmt"

	"golang.stackrox.io/kube-linter/internal/set"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/probeport/internal/params"
	"golang.stackrox.io/kube-linter/pkg/templates/util"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Probe Port",
		Key:         "probe-port",
		Description: "Flag unknown probe port",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(_ params.Params) (check.Func, error) {
			return util.PerNonInitContainerCheck(func(container *v1.Container) []diagnostic.Diagnostic {
				var portNames set.StringSet
				for _, port := range container.Ports {
					if name := port.Name; len(name) > 0 {
						portNames.Add(name)
					}
				}
				var results []diagnostic.Diagnostic
				for _, probe := range []*v1.Probe{container.LivenessProbe, container.ReadinessProbe, container.StartupProbe} {
					if probe == nil {
						continue
					}
					var port intstr.IntOrString
					if httpGet := probe.HTTPGet; httpGet != nil {
						port = httpGet.Port
					} else if tcpSocket := probe.TCPSocket; tcpSocket != nil {
						port = tcpSocket.Port
					} else {
						continue
					}
					if port.Type == intstr.String && !portNames.Contains(port.StrVal) && port.IntValue() == 0 {
						results = append(results, diagnostic.Diagnostic{
							Message: fmt.Sprintf("probe port %q does not match a port in container %q.", port.StrVal, container.Name),
						})
					}
				}
				return results
			}), nil
		}),
	})
}
