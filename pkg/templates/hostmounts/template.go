package hostmounts

import (
	"fmt"
	"path"
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

// matchesHostPath reports whether regex matches the host path as it is written
// in the manifest or in its cleaned form. The default dirs are exact matches
// such as "^/etc$", and a hostPath of "/etc/" mounts exactly the same directory
// without matching that pattern, so a stray separator would otherwise hide a
// sensitive mount. Cleaning also collapses "//etc" and "/etc/." for the same
// reason. The path as written is tried first so that a user-supplied pattern
// matching a literal trailing slash keeps working.
//
// path.Clean is used rather than filepath.Clean because a Kubernetes hostPath
// is always slash separated, whatever platform kube-linter runs on.
func matchesHostPath(regex *regexp.Regexp, hostPath string) bool {
	if regex.MatchString(hostPath) {
		return true
	}
	cleaned := path.Clean(hostPath)
	return cleaned != hostPath && regex.MatchString(cleaned)
}

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
						if !matchesHostPath(regex, v.HostPath.Path) {
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
