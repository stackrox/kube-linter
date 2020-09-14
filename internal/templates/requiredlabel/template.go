package requiredlabel

import (
	"regexp"

	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/diagnostic"
	"golang.stackrox.io/kube-linter/internal/extract"
	"golang.stackrox.io/kube-linter/internal/lintcontext"
	"golang.stackrox.io/kube-linter/internal/objectkinds"
	"golang.stackrox.io/kube-linter/internal/templates"
)

const (
	// TemplateName is the name of the required label template.
	TemplateName   = "required-label"
	keyParamName   = "key"
	valueParamName = "value"
)

func init() {
	templates.Register(check.Template{
		Name: TemplateName,
		SupportedObjectKinds: check.ObjectKindsDesc{
			ObjectKindNames: []string{objectkinds.Any},
		},
		Parameters: []check.ParameterDesc{
			{ParamName: keyParamName, Required: true, Description: "A regex for the key of the required label"},
			{ParamName: valueParamName, Required: false, Description: "A  regex for the value of the required label"},
		},
		Instantiate: func(params map[string]string) (check.Func, error) {
			keyRegexp, err := regexp.Compile(params[keyParamName])
			if err != nil {
				return nil, errors.Wrap(err, "invalid key")
			}
			var valueRegexp *regexp.Regexp
			if value := params[valueParamName]; value != "" {
				var err error
				valueRegexp, err = regexp.Compile(value)
				if err != nil {
					return nil, errors.Wrap(err, "invalid value")
				}
			}
			return func(_ *lintcontext.LintContext, object lintcontext.ObjectWithMetadata) []diagnostic.Diagnostic {
				labels := extract.Labels(object.K8sObject)
				for k, v := range labels {
					if !keyRegexp.MatchString(k) {
						continue
					}
					if valueRegexp != nil && !valueRegexp.MatchString(v) {
						continue
					}
					// Label found, nothing to return.
					return nil
				}
				return []diagnostic.Diagnostic{{
					Message: "no matching label found",
				}}
			}, nil
		},
	})
}
