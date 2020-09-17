package privileged

import (
	"fmt"

	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/diagnostic"
	"golang.stackrox.io/kube-linter/internal/extract"
	"golang.stackrox.io/kube-linter/internal/lintcontext"
	"golang.stackrox.io/kube-linter/internal/objectkinds"
	"golang.stackrox.io/kube-linter/internal/templates"
)

const (
	// TemplateName is the name of this template.
	TemplateName = "privileged"
)

func init() {
	templates.Register(check.Template{
		Name: TemplateName,
		SupportedObjectKinds: check.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters: nil,
		Instantiate: func(_ map[string]string) (check.Func, error) {
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
		},
	})
}
