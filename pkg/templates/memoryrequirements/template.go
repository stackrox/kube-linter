package memoryrequirements

import (
	"fmt"

	"golang.stackrox.io/kube-linter/internal/pointers"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/memoryrequirements/internal/params"
	"golang.stackrox.io/kube-linter/pkg/templates/util"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	bytesInMB = 1024 * 1024
)

func process(results *[]diagnostic.Diagnostic, containerName, requirementsType string, quantity *resource.Quantity, lowerBoundBytes int, upperBoundBytes *int) {
	if util.ValueInRange(int(quantity.Value()), lowerBoundBytes, upperBoundBytes) {
		*results = append(*results, diagnostic.Diagnostic{
			Message: fmt.Sprintf("container %q has memory %s %s", containerName, requirementsType, quantity),
		})
	}
}

func init() {
	templates.Register(check.Template{
		HumanName:   "Memory Requirements",
		Key:         "memory-requirements",
		Description: "Flag containers with memory requirements in the given range",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			lowerBoundBytes := p.LowerBoundMB * bytesInMB
			var upperBoundBytes *int
			if p.UpperBoundMB != nil {
				upperBoundBytes = pointers.Int((*p.UpperBoundMB) * bytesInMB)
			}
			return util.PerContainerCheck(func(container *v1.Container) []diagnostic.Diagnostic {
				var results []diagnostic.Diagnostic
				if p.RequirementsType == "request" || p.RequirementsType == "any" {
					process(&results, container.Name, "request", container.Resources.Requests.Memory(), lowerBoundBytes, upperBoundBytes)
				}
				if p.RequirementsType == "limit" || p.RequirementsType == "any" {
					process(&results, container.Name, "limit", container.Resources.Limits.Memory(), lowerBoundBytes, upperBoundBytes)
				}
				return results
			}), nil
		}),
	})
}
