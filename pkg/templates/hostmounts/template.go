package hostmounts

import (
	"fmt"
	"regexp"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/hostmounts/internal/params"
)

const (
	templateKey = "host-mounts"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Host Mounts",
		Key:         templateKey,
		Description: "Flag volume mounts of sensitive system directories",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			compiledRegexes := make([]*regexp.Regexp, 0, len(p.Dirs))
			for _, dir := range p.Dirs {
				r, err := regexp.Compile(dir)
				if err != nil {
					return nil, fmt.Errorf("invalid regex %s: %w", dir, err)
				}
				compiledRegexes = append(compiledRegexes, r)
			}
			return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				podSpec, found := extract.PodSpec(object.K8sObject)
				if !found {
					return nil
				}
				var results []diagnostic.Diagnostic
				containers := podSpec.AllContainers()
				for _, v := range podSpec.Volumes {
					if v.HostPath == nil {
						continue
					}
					for _, regex := range compiledRegexes {
						if !regex.MatchString(v.HostPath.Path) {
							continue
						}
						for _, container := range containers {
							for _, mount := range container.VolumeMounts {
								if mount.Name == v.Name {
									results = append(results, diagnostic.Diagnostic{
										Message: fmt.Sprintf("host system directory %q is mounted on container %q", v.HostPath.Path, container.Name)})
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
