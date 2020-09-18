package runasnonroot

import (
	"fmt"

	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/diagnostic"
	"golang.stackrox.io/kube-linter/internal/extract"
	"golang.stackrox.io/kube-linter/internal/lintcontext"
	"golang.stackrox.io/kube-linter/internal/objectkinds"
	"golang.stackrox.io/kube-linter/internal/templates"
)

func init() {
	templates.Register(check.Template{
		Name:        "run-as-non-root",
		Description: "Flag containers without runAsUser specified",
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
				var podSCValue bool
				if podSpec.SecurityContext != nil && podSpec.SecurityContext.RunAsNonRoot != nil {
					podSCValue = *podSpec.SecurityContext.RunAsNonRoot
				}
				for i := range podSpec.Containers {
					container := podSpec.Containers[i]
					var runAsNonRoot *bool
					if container.SecurityContext != nil {
						runAsNonRoot = container.SecurityContext.RunAsNonRoot
					}
					if runAsNonRoot != nil && *runAsNonRoot {
						continue
					}
					if runAsNonRoot == nil && podSCValue {
						continue
					}
					results = append(results, diagnostic.Diagnostic{Message: fmt.Sprintf("container %q is not set to runAsNonRoot", container.Name)})
				}
				return results
			}, nil
		},
	})
}
