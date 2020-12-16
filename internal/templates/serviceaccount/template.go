package serviceaccount

import (
	"fmt"

	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/diagnostic"
	"golang.stackrox.io/kube-linter/internal/extract"
	"golang.stackrox.io/kube-linter/internal/lintcontext"
	"golang.stackrox.io/kube-linter/internal/matcher"
	"golang.stackrox.io/kube-linter/internal/objectkinds"
	"golang.stackrox.io/kube-linter/internal/stringutils"
	"golang.stackrox.io/kube-linter/internal/templates"
	"golang.stackrox.io/kube-linter/internal/templates/serviceaccount/internal/params"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Service Account",
		Key:         "service-account",
		Description: "Flag containers which use a matching service account",
		SupportedObjectKinds: check.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			saMatcher, err := matcher.ForString(p.ServiceAccount)
			if err != nil {
				return nil, err
			}
			return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				podSpec, found := extract.PodSpec(object.K8sObject)
				if !found {
					return nil
				}
				sa := stringutils.OrDefault(podSpec.ServiceAccountName, podSpec.DeprecatedServiceAccount)
				if saMatcher(sa) {
					return []diagnostic.Diagnostic{{Message: fmt.Sprintf("found matching serviceAccount (%q)", sa)}}
				}
				return nil
			}, nil
		}),
	})
}
