package writablehostmount

import (
	"fmt"

	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/diagnostic"
	"golang.stackrox.io/kube-linter/internal/extract"
	"golang.stackrox.io/kube-linter/internal/lintcontext"
	"golang.stackrox.io/kube-linter/internal/objectkinds"
	"golang.stackrox.io/kube-linter/internal/templates"
	"golang.stackrox.io/kube-linter/internal/templates/writablehostmount/internal/params"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Writable Host Mounts",
		Key:         "writable-host-mount",
		Description: "Flag containers that have mounted a directory on the host as writable",
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
				hostPaths := make(map[string]string)
				for _, volume := range podSpec.Volumes {
					if volume.HostPath != nil {
						hostPaths[volume.Name] = volume.HostPath.Path
					}
				}
				if len(hostPaths) == 0 {
					return nil
				}
				var results []diagnostic.Diagnostic
				for _, container := range podSpec.Containers {
					for _, mount := range container.VolumeMounts {
						if mount.ReadOnly {
							continue
						}
						if hostPath, exists := hostPaths[mount.Name]; exists {
							results = append(results, diagnostic.Diagnostic{Message: fmt.Sprintf("container %s mounts path %s on the host as writable", container.Name, hostPath)})
						}
					}
				}
				return results
			}, nil
		}),
	})
}
