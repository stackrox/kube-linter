package runasnonroot

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/runasnonroot/internal/params"
	v1 "k8s.io/api/core/v1"
)

const templateKey = "run-as-non-root"

func effectiveRunAsNonRoot(podSC *v1.PodSecurityContext, containerSC *v1.SecurityContext) bool {
	if containerSC != nil && containerSC.RunAsNonRoot != nil {
		return *containerSC.RunAsNonRoot
	}
	if podSC != nil && podSC.RunAsNonRoot != nil {
		return *podSC.RunAsNonRoot
	}
	return false
}

func effectiveRunAsUser(podSC *v1.PodSecurityContext, containerSC *v1.SecurityContext) *int64 {
	if containerSC != nil && containerSC.RunAsUser != nil {
		return containerSC.RunAsUser
	}
	if podSC != nil {
		return podSC.RunAsUser
	}
	return nil
}

func effectiveRunAsGroup(podSC *v1.PodSecurityContext, containerSC *v1.SecurityContext) *int64 {
	if containerSC != nil && containerSC.RunAsGroup != nil {
		return containerSC.RunAsGroup
	}
	if podSC != nil {
		return podSC.RunAsGroup
	}
	return nil
}

func isNonZero(number *int64) bool {
	return number != nil && *number > 0
}

func init() {
	templates.Register(check.Template{
		HumanName:   "Run as non-root",
		Key:         templateKey,
		Description: "Flag containers set to run as a root user or group",
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
				var results []diagnostic.Diagnostic
				for _, container := range podSpec.AllContainers() {
					runAsUser := effectiveRunAsUser(podSpec.SecurityContext, container.SecurityContext)
					runAsGroup := effectiveRunAsGroup(podSpec.SecurityContext, container.SecurityContext)
					// runAsUser and runAsGroup explicitly set to non-root. All good.
					if isNonZero(runAsUser) && isNonZero(runAsGroup) {
						continue
					}

					// runAsGroup set to 0
					if runAsGroup != nil && *runAsGroup == 0 {
						results = append(results, diagnostic.Diagnostic{
							Message: fmt.Sprintf("container %q has runAsGroup set to %d", container.Name, *runAsGroup),
						})
					}
					// runAsGroup is not set.
					if runAsGroup == nil {
						results = append(results, diagnostic.Diagnostic{
							Message: fmt.Sprintf("container %q does not have runAsGroup set", container.Name),
						})
					}

					runAsNonRoot := effectiveRunAsNonRoot(podSpec.SecurityContext, container.SecurityContext)
					if runAsNonRoot {
						// runAsNonRoot set, but runAsUser set to 0. This will result in a runtime failure.
						if runAsUser != nil && *runAsUser == 0 {
							results = append(results, diagnostic.Diagnostic{
								Message: fmt.Sprintf("container %q is set to runAsNonRoot, but runAsUser set to %d", container.Name, *runAsUser),
							})
						}
						continue
					}
					results = append(results, diagnostic.Diagnostic{Message: fmt.Sprintf("container %q is not set to runAsNonRoot", container.Name)})
				}
				return results
			}, nil
		}),
	})
}
