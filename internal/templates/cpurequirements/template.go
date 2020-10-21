package cpurequirements

import (
	"fmt"

	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/diagnostic"
	"golang.stackrox.io/kube-linter/internal/extract"
	"golang.stackrox.io/kube-linter/internal/lintcontext"
	"golang.stackrox.io/kube-linter/internal/objectkinds"
	"golang.stackrox.io/kube-linter/internal/templates"
	"golang.stackrox.io/kube-linter/internal/templates/cpurequirements/internal/params"
	"golang.stackrox.io/kube-linter/internal/templates/util"
	"k8s.io/apimachinery/pkg/api/resource"
)

func process(results *[]diagnostic.Diagnostic, containerName, requirementsType string, quantity *resource.Quantity, lowerBound int, upperBound *int) {
	if util.ValueInRange(int(quantity.MilliValue()), lowerBound, upperBound) {
		*results = append(*results, diagnostic.Diagnostic{
			Message: fmt.Sprintf("container %q has cpu %s %s", containerName, requirementsType, quantity),
		})
	}

}

func init() {
	templates.Register(check.Template{
		HumanName:   "CPU Requirements",
		Key:         "cpu-requirements",
		Description: "Flag containers with CPU requirements in the given range",
		SupportedObjectKinds: check.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			return func(_ *lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				podSpec, found := extract.PodSpec(object.K8sObject)
				if !found {
					return nil
				}

				var results []diagnostic.Diagnostic
				for _, container := range podSpec.Containers {
					if p.RequirementsType == "request" || p.RequirementsType == "any" {
						process(&results, container.Name, "request", container.Resources.Requests.Cpu(), p.LowerBoundMillis, p.UpperBoundMillis)
					}
					if p.RequirementsType == "limit" || p.RequirementsType == "any" {
						process(&results, container.Name, "limit", container.Resources.Limits.Cpu(), p.LowerBoundMillis, p.UpperBoundMillis)
					}
				}
				return results
			}, nil
		}),
	})
}
