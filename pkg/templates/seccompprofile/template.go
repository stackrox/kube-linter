package seccompprofile

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/seccompprofile/internal/params"
	v1 "k8s.io/api/core/v1"
)

const templateKey = "seccomp-profile"

// effectiveSeccompProfile returns the seccomp profile that applies to a
// container. A container-level profile takes precedence over the pod-level
// one, mirroring how Kubernetes resolves the securityContext fields.
func effectiveSeccompProfile(podSC *v1.PodSecurityContext, containerSC *v1.SecurityContext) *v1.SeccompProfile {
	if containerSC != nil && containerSC.SeccompProfile != nil {
		return containerSC.SeccompProfile
	}
	if podSC != nil {
		return podSC.SeccompProfile
	}
	return nil
}

// hasSeccompProfile returns true if the given profile actually confines the
// container. An unset profile, an empty profile type, and an explicit
// "Unconfined" profile all leave the container without seccomp filtering.
func hasSeccompProfile(profile *v1.SeccompProfile) bool {
	return profile != nil &&
		profile.Type != "" &&
		profile.Type != v1.SeccompProfileTypeUnconfined
}

func init() {
	templates.Register(check.Template{
		HumanName:   "Seccomp Profile",
		Key:         templateKey,
		Description: "Flag containers without a seccomp profile",
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
					if hasSeccompProfile(effectiveSeccompProfile(podSpec.SecurityContext, container.SecurityContext)) {
						continue
					}
					results = append(results, diagnostic.Diagnostic{Message: fmt.Sprintf(
						"container %q does not have a seccomp profile set. "+
							"Set seccompProfile.type to %q or %q in the securityContext of the container or the pod.",
						container.Name, v1.SeccompProfileTypeRuntimeDefault, v1.SeccompProfileTypeLocalhost)})
				}
				return results
			}, nil
		}),
	})
}
