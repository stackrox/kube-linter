package priorityclassname

import (
	"fmt"
	"strings"

	"slices"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/priorityclassname/internal/params"
)

const (
	templateKey = "priority-class-name"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Priority class name",
		Key:         templateKey,
		Description: "Flag applications running with invalid priority class name.",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				spec, found := extract.PodSpec(object.K8sObject)
				isEmpty := strings.TrimSpace(spec.PriorityClassName) == ""
				isAccepted := slices.Contains(p.AcceptedPriorityClassNames, spec.PriorityClassName)
				if !found || isEmpty || isAccepted {
					return nil
				}
				return []diagnostic.Diagnostic{
					{Message: fmt.Sprintf("object has a priority class name defined with '%s' but the only accepted priority class names are '%s'", spec.PriorityClassName, p.AcceptedPriorityClassNames)},
				}
			}, nil
		}),
	})
}
