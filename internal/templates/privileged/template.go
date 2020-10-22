package privileged

import (
	"fmt"

	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/diagnostic"
	"golang.stackrox.io/kube-linter/internal/extract"
	"golang.stackrox.io/kube-linter/internal/lintcontext"
	"golang.stackrox.io/kube-linter/internal/objectkinds"
	"golang.stackrox.io/kube-linter/internal/templates"
	"golang.stackrox.io/kube-linter/internal/templates/privileged/internal/params"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Privileged Containers",
		Key:         "privileged",
		Description: "Flag privileged containers",
		SupportedObjectKinds: check.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(_ params.Params) (check.Func, error) {
			return func(_ *lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				podSpec, found := extract.PodSpec(object.K8sObject)
				if !found {
					return nil
				}
				var results []diagnostic.Diagnostic
				for i := range podSpec.Containers {
					container := &podSpec.Containers[i]
					if securityContext := container.SecurityContext; securityContext != nil {
						if securityContext.Privileged != nil && *securityContext.Privileged {
							results = append(results, diagnostic.Diagnostic{Message: fmt.Sprintf("container %q is privileged", container.Name)})
						}
					}
				}
				return results
			}, nil
		}),
	})
}
