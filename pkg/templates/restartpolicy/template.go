package restartpolicy

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/restartpolicy/internal/params"
	coreV1 "k8s.io/api/core/v1"
)

const (
	templateKey = "restart-policy"
)

var acceptedRestartPolicies = []coreV1.RestartPolicy{coreV1.RestartPolicyAlways, coreV1.RestartPolicyOnFailure}

func init() {
	templates.Register(check.Template{
		HumanName:   "Restart policy",
		Key:         templateKey,
		Description: "Flag applications running without the restart policy.",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				spec, found := extract.PodSpec(object.K8sObject)
				if !found {
					return nil
				}
				for _, policy := range acceptedRestartPolicies {
					if spec.RestartPolicy == policy {
						return nil
					}
				}
				return []diagnostic.Diagnostic{
					{Message: fmt.Sprintf("object has a restart policy defined with '%s' but the only accepted restart policies are '%s'", spec.RestartPolicy, acceptedRestartPolicies)},
				}
			}, nil
		}),
	})
}
