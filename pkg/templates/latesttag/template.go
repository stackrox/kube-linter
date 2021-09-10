package latesttag

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/latesttag/internal/params"
)

const (
	templateKey = "latest-tag"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Latest Tag",
		Key:         templateKey,
		Description: "Flag applications running container images that do not satisfies \"allowList\" & \"blockList\" parameters criteria.",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {

			blockedRegexes := make([]*regexp.Regexp, 0, len(p.BlockList))
			for _, res := range p.BlockList {
				rg, err := regexp.Compile(res)
				if err != nil {
					return nil, errors.Wrapf(err, "invalid regex %s", res)
				}
				blockedRegexes = append(blockedRegexes, rg)
			}

			allowedRegexes := make([]*regexp.Regexp, 0, len(p.AllowList))
			for _, res := range p.AllowList {
				rg, err := regexp.Compile(res)
				if err != nil {
					return nil, errors.Wrapf(err, "invalid regex %s", res)
				}
				allowedRegexes = append(allowedRegexes, rg)
			}

			if len(blockedRegexes) > 0 && len(allowedRegexes) > 0 {
				err := fmt.Errorf("check has both \"allowList\" & \"blockList\" parameter's values set")
				return nil, errors.Wrapf(err, "Only one of the paramater lists can be used at a time.")
			}

			return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				podSpec, found := extract.PodSpec(object.K8sObject)
				if !found {
					return nil
				}

				var results []diagnostic.Diagnostic

				for _, container := range podSpec.Containers {
					if len(blockedRegexes) > 0 && isInList(blockedRegexes, container.Image) {
						results = append(results, diagnostic.Diagnostic{Message: fmt.Sprintf("The container %q is using an invalid container image, %q. Please use images that are not blocked by the `BlockList` criteria : %q", container.Name, container.Image, blockedRegexes)})
					} else if len(allowedRegexes) > 0 && !isInList(allowedRegexes, container.Image) {
						results = append(results, diagnostic.Diagnostic{Message: fmt.Sprintf("The container %q is using an invalid container image, %q. Please use images that satisfies the `AllowList` criteria : %q", container.Name, container.Image, allowedRegexes)})
					}
				}
				return results
			}, nil
		}),
	})
}

// isInList returns true if a match found in the list for the given name
func isInList(regexlist []*regexp.Regexp, name string) bool {
	for _, regex := range regexlist {
		if regex.MatchString(name) {
			return true
		}
	}
	return false
}
