package imagedigest

import (
	"fmt"
	"regexp"

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
	templateKey = "image-digest"
)

// digestRegex matches the digest portion of an image reference, for example
// "@sha256:2cf0d5dcdb7680af443d09f850490eb79d9245b6e5c9b04b9a70ea38513b6c00".
// A reference carrying a digest pins the exact image content, unlike a tag,
// which is mutable and can be re-pointed at a different image at any time.
var digestRegex = regexp.MustCompile(`@[a-zA-Z0-9][a-zA-Z0-9._-]*:[a-fA-F0-9]{32,}`)

func init() {
	templates.Register(check.Template{
		HumanName:   "Image Digest",
		Key:         templateKey,
		Description: "Flag containers whose image is not pinned to an immutable digest.",
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
				if container.Image == "" || hasDigest(container.Image) || isAllowed(allowedRegexes, container.Image) {
					return nil
				}
				return []diagnostic.Diagnostic{{Message: fmt.Sprintf(
					"container %q is using image %q which is not pinned to a digest. "+
						"Pin the image to an immutable digest (for example, \"image@sha256:<digest>\") so that a tag cannot be re-pointed at a different image.",
					container.Name, container.Image)}}
			}), nil
		}),
	})
}

// hasDigest returns true if the image reference carries a content digest.
func hasDigest(image string) bool {
	return digestRegex.MatchString(image)
}

// isAllowed returns true if the image matches any of the exemption regexes.
func isAllowed(regexes []*regexp.Regexp, image string) bool {
	for _, rg := range regexes {
		if rg.MatchString(image) {
			return true
		}
	}
	return false
}
