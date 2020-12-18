package deprecatedserviceaccount

import (
	"fmt"

	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/diagnostic"
	"golang.stackrox.io/kube-linter/internal/extract"
	"golang.stackrox.io/kube-linter/internal/lintcontext"
	"golang.stackrox.io/kube-linter/internal/objectkinds"
	"golang.stackrox.io/kube-linter/internal/templates"
	"golang.stackrox.io/kube-linter/internal/templates/deprecatedserviceaccount/internal/params"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Deprecated Service Account Field",
		Key:         "deprecated-service-account-field",
		Description: "Flag uses of the deprecated serviceAccount field, which should be migrated to serviceAccountName",
		SupportedObjectKinds: check.ObjectKindsDesc{
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
				if sa := podSpec.DeprecatedServiceAccount; sa != "" {
					return []diagnostic.Diagnostic{{Message: fmt.Sprintf(
						"serviceAccount is specified (%s), but this field is deprecated; use serviceAccountName instead", sa)}}
				}
				return nil
			}, nil
		}),
	})
}
