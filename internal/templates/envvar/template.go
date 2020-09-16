package envvar

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/diagnostic"
	"golang.stackrox.io/kube-linter/internal/extract"
	"golang.stackrox.io/kube-linter/internal/lintcontext"
	"golang.stackrox.io/kube-linter/internal/matcher"
	"golang.stackrox.io/kube-linter/internal/objectkinds"
	"golang.stackrox.io/kube-linter/internal/templates"
)

const (
	// TemplateName is the name of this template.
	TemplateName   = "env-var"
	nameParamName  = "name"
	valueParamName = "value"
)

func init() {
	templates.Register(check.Template{
		Name: TemplateName,
		SupportedObjectKinds: check.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters: []check.ParameterDesc{
			{ParamName: nameParamName, Required: true, Description: "A regex for the env var name"},
			{ParamName: valueParamName, Description: "A regex for the env var value"},
		},
		Instantiate: func(params map[string]string) (check.Func, error) {
			nameMatcher, err := matcher.ForString(params[nameParamName])
			if err != nil {
				return nil, errors.Wrap(err, "invalid key")
			}
			valueMatcher, err := matcher.ForString(params[valueParamName])
			if err != nil {
				return nil, errors.Wrap(err, "invalid value")
			}
			return func(_ *lintcontext.LintContext, object lintcontext.ObjectWithMetadata) []diagnostic.Diagnostic {
				podSpec, found := extract.PodSpec(object.K8sObject)
				if !found {
					return nil
				}
				var results []diagnostic.Diagnostic
				for i := range podSpec.Containers {
					container := &podSpec.Containers[i]
					for _, envVar := range container.Env {
						if nameMatcher(envVar.Name) && valueMatcher(envVar.Value) {
							results = append(results, diagnostic.Diagnostic{
								Message: fmt.Sprintf("environment variable %s in container %q found", envVar.Name, container.Name),
							})
						}
					}
				}
				return results
			}, nil
		},
	})
}
