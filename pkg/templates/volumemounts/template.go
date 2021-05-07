package volumemounts

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/volumemounts/internal/params"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Volume Mounts",
		Key:         "volume-mounts",
		Description: "Flag volume mounts of sensitive system directories",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				podSpec, found := extract.PodSpec(object.K8sObject)
				if !found {
					return nil
				}
				var results []diagnostic.Diagnostic
				containers := podSpec.Containers
				for _, v := range podSpec.Volumes {
					if v.HostPath != nil {
						for _, dir := range p.SensitiveSysDirs {
							if v.HostPath.Path == dir {
								for _, container := range containers {
									for _, mount := range container.VolumeMounts {
										if mount.Name == v.Name {
											results = append(results, diagnostic.Diagnostic{
												Message: fmt.Sprintf("sensitive host system directory %s is mounted on container %s", dir, container.Name)})
										}
									}
								}
							}
						}
					}
				}
				return results
			}, nil
		}),
	})
}
