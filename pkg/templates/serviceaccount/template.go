package serviceaccount

import (
	"fmt"

	"golang.stackrox.io/kube-linter/internal/stringutils"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/matcher"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/serviceaccount/internal/params"
)

const (
	templateKey = "service-account"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Service Account",
		Key:         templateKey,
		Description: "Flag containers which use a matching service account",
		SupportedObjectKinds: config.ObjectKindsDesc{
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
				if podSpec.AutomountServiceAccountToken != nil && !*podSpec.AutomountServiceAccountToken {
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
