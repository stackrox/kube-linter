package ports

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/matcher"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/ports/internal/params"
	"golang.stackrox.io/kube-linter/pkg/templates/util"
	v1 "k8s.io/api/core/v1"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Ports",
		Key:         "ports",
		Description: "Flag containers exposing ports under protocols that match the supplied parameters",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			protocolMatcher, err := matcher.ForString(p.Protocol)
			if err != nil {
				return nil, fmt.Errorf("invalid protocol: %w", err)
			}
			return util.PerContainerCheck(func(container *v1.Container) []diagnostic.Diagnostic {
				var results []diagnostic.Diagnostic
				for _, port := range container.Ports {
					// The k8s protocol defaults to TCP even if not set in the YAML.
					protocol := string(port.Protocol)
					if protocol == "" {
						protocol = "TCP"
					}
					if int(port.ContainerPort) == p.Port && protocolMatcher(protocol) {
						results = append(results, diagnostic.Diagnostic{
							Message: fmt.Sprintf("port %d and protocol %s in container %q found",
								port.ContainerPort, protocol, container.Name),
						})
					}
				}
				return results
			}), nil
		}),
	})
}
