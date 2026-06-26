package imagedigest

import (
	"fmt"
	"regexp"
	"strings"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/imagedigest/internal/params"
	"golang.stackrox.io/kube-linter/pkg/templates/util"
	v1 "k8s.io/api/core/v1"
)

const (
	templateKey = "image-not-pinned-by-digest"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Image Not Pinned By Digest",
		Key:         templateKey,
		Description: "Flag container images that are not pinned to a specific digest",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			allowedRegexes := make([]*regexp.Regexp, 0, len(p.AllowList))
			for _, res := range p.AllowList {
				rg, err := regexp.Compile(res)
				if err != nil {
					return nil, fmt.Errorf("invalid regex %s: %w", res, err)
				}
				allowedRegexes = append(allowedRegexes, rg)
			}

			return util.PerContainerCheck(func(container *v1.Container) (results []diagnostic.Diagnostic) {
				if isAllowed(allowedRegexes, container.Image) {
					return nil
				}
				if !strings.Contains(container.Image, "@sha256:") {
					results = append(results, diagnostic.Diagnostic{
						Message: fmt.Sprintf("container %q image %q is not pinned by digest", container.Name, container.Image),
					})
				}
				return results
			}), nil
		}),
	})
}

func isAllowed(regexlist []*regexp.Regexp, name string) bool {
	for _, regex := range regexlist {
		if regex.MatchString(name) {
			return true
		}
	}
	return false
}
