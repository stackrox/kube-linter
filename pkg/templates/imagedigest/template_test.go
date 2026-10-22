package imagedigest

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/imagedigest/internal/params"

	v1 "k8s.io/api/core/v1"
)

var (
	containerName = "test-container"
)

func TestImageDigest(t *testing.T) {
	suite.Run(t, new(ImageDigestTestSuite))
}

type ImageDigestTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *ImageDigestTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *ImageDigestTestSuite) addDeploymentWithContainerImage(name, containerImage string) {
	s.ctx.AddMockDeployment(s.T(), name)
	s.ctx.AddContainerToDeployment(s.T(), name, v1.Container{Name: containerName, Image: containerImage})
}

func (s *ImageDigestTestSuite) TestUnpinnedImagesFlagged() {
	const (
		taggedImage   = "dep-with-tagged-image"
		latestImage   = "dep-with-latest-image"
		untaggedImage = "dep-with-untagged-image"
	)

	s.addDeploymentWithContainerImage(taggedImage, "example.com/app:v1.2.3")
	s.addDeploymentWithContainerImage(latestImage, "example.com/app:latest")
	s.addDeploymentWithContainerImage(untaggedImage, "example.com/app")

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				taggedImage: {
					{Message: "container \"test-container\" is using image \"example.com/app:v1.2.3\" which is not pinned to a digest. Pin the image to an immutable digest (for example, \"image@sha256:<digest>\") so that a tag cannot be re-pointed at a different image."},
				},
				latestImage: {
					{Message: "container \"test-container\" is using image \"example.com/app:latest\" which is not pinned to a digest. Pin the image to an immutable digest (for example, \"image@sha256:<digest>\") so that a tag cannot be re-pointed at a different image."},
				},
				untaggedImage: {
					{Message: "container \"test-container\" is using image \"example.com/app\" which is not pinned to a digest. Pin the image to an immutable digest (for example, \"image@sha256:<digest>\") so that a tag cannot be re-pointed at a different image."},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *ImageDigestTestSuite) TestDigestPinnedImagesPass() {
	const (
		digestOnly    = "dep-with-digest-only"
		tagAndDigest  = "dep-with-tag-and-digest"
		allowedByList = "dep-with-allowed-image"
	)

	s.addDeploymentWithContainerImage(digestOnly, "example.com/app@sha256:2cf0d5dcdb7680af443d09f850490eb79d9245b6e5c9b04b9a70ea38513b6c00")
	s.addDeploymentWithContainerImage(tagAndDigest, "example.com/app:v1.2.3@sha256:2cf0d5dcdb7680af443d09f850490eb79d9245b6e5c9b04b9a70ea38513b6c00")
	s.addDeploymentWithContainerImage(allowedByList, "localhost:5000/dev-image")

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				AllowList: []string{"^localhost:5000/"},
			},
			Diagnostics:              map[string][]diagnostic.Diagnostic{},
			ExpectInstantiationError: false,
		},
	})
}

func (s *ImageDigestTestSuite) TestInvalidAllowListRegex() {
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				AllowList: []string{"["},
			},
			Diagnostics:              map[string][]diagnostic.Diagnostic{},
			ExpectInstantiationError: true,
		},
	})
}
