package deprecatedserviceaccount

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/deprecatedserviceaccount/internal/params"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Deprecated Service Account Field",
		Key:         "deprecated-service-account-field",
		Description: "Flag uses of the deprecated serviceAccount field, which should be migrated to serviceAccountName",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(_ params.Params) (check.Func, error) {
			return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				podSpec, found := extract.PodSpec(object.K8sObject)
				if !found {
					return nil
				}

				sa := podSpec.DeprecatedServiceAccount
				san := podSpec.ServiceAccountName
				if sa != "" && sa != san {
					if san == "" { // only serviceAccount is specified
						return []diagnostic.Diagnostic{{Message: fmt.Sprintf(
							"serviceAccount is specified (%s), but this field is deprecated; use serviceAccountName instead", sa)}}
					}
					// serviceAccount and serviceAccountName both specified but do not match
					return []diagnostic.Diagnostic{{Message: fmt.Sprintf(
						"serviceAccount (%s) and serviceAccountName (%s) are both specified with non-matching values. serviceAccount is deprecated; unspecify serviceAccount or make values match.", sa, san)}}
				}
				return nil
			}, nil
		}),
	})
}
