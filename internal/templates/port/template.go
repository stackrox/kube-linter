package port

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/internal/matcher"

	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/diagnostic"
	"golang.stackrox.io/kube-linter/internal/objectkinds"
	"golang.stackrox.io/kube-linter/internal/templates"
	"golang.stackrox.io/kube-linter/internal/templates/port/internal/params"
	"golang.stackrox.io/kube-linter/internal/templates/util"
	v1 "k8s.io/api/core/v1"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Ports",
		Key:         "ports",
		Description: "Flag containers exposing ports under protocols that match the supplied parameters",
		SupportedObjectKinds: check.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			protocolMatcher, err := matcher.ForString(p.Protocol)
			if err != nil {
				return nil, errors.Wrap(err, "invalid protocol")
			}
			return util.PerContainerCheck(func(container *v1.Container) []diagnostic.Diagnostic {
				var results []diagnostic.Diagnostic
				for _, port := range container.Ports {
					if int(port.ContainerPort) == p.Port && protocolMatcher(string(port.Protocol)) {
						results = append(results, diagnostic.Diagnostic{
							Message: fmt.Sprintf("port %d and protocol %s in container %q found", port.ContainerPort, string(port.Protocol), container.Name),
						})
					}
				}
				return results
			}), nil
		}),
	})
}
