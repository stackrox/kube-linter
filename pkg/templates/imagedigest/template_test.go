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

func (s *ImageDigestTestSuite) addDeploymentWithInitContainerImage(name, containerImage string) {
	s.ctx.AddMockDeployment(s.T(), name)
	s.ctx.AddInitContainerToDeployment(s.T(), name, v1.Container{Name: "test-init-container", Image: containerImage})
}

func (s *ImageDigestTestSuite) addDeploymentWithEphemeralContainerImage(name, containerImage string) {
	s.ctx.AddMockDeployment(s.T(), name)
	s.ctx.AddEphemeralContainerToDeployment(s.T(), name, v1.EphemeralContainer{
		EphemeralContainerCommon: v1.EphemeralContainerCommon{Name: "test-ephemeral-container", Image: containerImage},
	})
}

func (s *ImageDigestTestSuite) TestUnpinnedImagesFlagged() {
	const (
		taggedImage       = "dep-with-tagged-image"
		latestImage       = "dep-with-latest-image"
		untaggedImage     = "dep-with-untagged-image"
		shortDigestImg    = "dep-with-short-sha256-digest"
		shortSha512Img    = "dep-with-short-sha512-digest"
		malformedDigest   = "dep-with-malformed-digest"
		initUnpinned      = "dep-with-unpinned-init-container"
		ephemeralUnpinned = "dep-with-unpinned-ephemeral-container"
	)

	s.addDeploymentWithContainerImage(taggedImage, "example.com/app:v1.2.3")
	s.addDeploymentWithContainerImage(latestImage, "example.com/app:latest")
	s.addDeploymentWithContainerImage(untaggedImage, "example.com/app")
	s.addDeploymentWithContainerImage(shortDigestImg, "example.com/app@sha256:2cf0d5dcdb7680af443d09f850490eb79d9245b6e")
	s.addDeploymentWithContainerImage(shortSha512Img, "example.com/app@sha512:401b09eab3c013d4ca54922bb802bec8")
	s.addDeploymentWithContainerImage(malformedDigest, "example.com/app@sha256:2cf0d5dcdb7680af443d09f850490eb79d9245b6e5c9b04b9a70ea38513b6c00:extra")
	s.addDeploymentWithInitContainerImage(initUnpinned, "example.com/init:v1")
	s.addDeploymentWithEphemeralContainerImage(ephemeralUnpinned, "example.com/ephemeral:v1")

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
				shortDigestImg: {
					{Message: "container \"test-container\" is using image \"example.com/app@sha256:2cf0d5dcdb7680af443d09f850490eb79d9245b6e\" which is not pinned to a digest. Pin the image to an immutable digest (for example, \"image@sha256:<digest>\") so that a tag cannot be re-pointed at a different image."},
				},
				malformedDigest: {
					{Message: "container \"test-container\" is using image \"example.com/app@sha256:2cf0d5dcdb7680af443d09f850490eb79d9245b6e5c9b04b9a70ea38513b6c00:extra\" which is not pinned to a digest. Pin the image to an immutable digest (for example, \"image@sha256:<digest>\") so that a tag cannot be re-pointed at a different image."},
				},
				shortSha512Img: {
					{Message: "container \"test-container\" is using image \"example.com/app@sha512:401b09eab3c013d4ca54922bb802bec8\" which is not pinned to a digest. Pin the image to an immutable digest (for example, \"image@sha256:<digest>\") so that a tag cannot be re-pointed at a different image."},
				},
				initUnpinned: {
					{Message: "container \"test-init-container\" is using image \"example.com/init:v1\" which is not pinned to a digest. Pin the image to an immutable digest (for example, \"image@sha256:<digest>\") so that a tag cannot be re-pointed at a different image."},
				},
				ephemeralUnpinned: {
					{Message: "container \"test-ephemeral-container\" is using image \"example.com/ephemeral:v1\" which is not pinned to a digest. Pin the image to an immutable digest (for example, \"image@sha256:<digest>\") so that a tag cannot be re-pointed at a different image."},
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
		sha512Digest  = "dep-with-sha512-digest"
		allowedByList = "dep-with-allowed-image"
	)

	s.addDeploymentWithContainerImage(digestOnly, "example.com/app@sha256:2cf0d5dcdb7680af443d09f850490eb79d9245b6e5c9b04b9a70ea38513b6c00")
	s.addDeploymentWithContainerImage(tagAndDigest, "example.com/app:v1.2.3@sha256:2cf0d5dcdb7680af443d09f850490eb79d9245b6e5c9b04b9a70ea38513b6c00")
	s.addDeploymentWithContainerImage(sha512Digest, "example.com/app@sha512:401b09eab3c013d4ca54922bb802bec8fd5318192b0a75f201d8b372742401b09eab3c013d4ca54922bb802bec8fd5318192b0a75f201d8b3727420123456789")
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

func (s *ImageDigestTestSuite) TestOtherAlgorithmDigests() {
	const (
		sha384Digest     = "dep-with-sha384-digest"
		shortSha384      = "dep-with-short-sha384-digest"
		otherAlgo32Chars = "dep-with-32-char-other-digest"
		emptyImage       = "dep-with-empty-image"
	)

	// Algorithms other than sha256 and sha512 count as a digest when the hex
	// part is at least 32 characters, so a full sha384 digest and a 32-char
	// digest pass, while a 16-char one is flagged.
	s.addDeploymentWithContainerImage(sha384Digest, "example.com/app@sha384:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	s.addDeploymentWithContainerImage(shortSha384, "example.com/app@sha384:0123456789abcdef")
	s.addDeploymentWithContainerImage(otherAlgo32Chars, "example.com/app@md5:0123456789abcdef0123456789abcdef")
	s.addDeploymentWithContainerImage(emptyImage, "")

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				shortSha384: {
					{Message: "container \"test-container\" is using image \"example.com/app@sha384:0123456789abcdef\" which is not pinned to a digest. Pin the image to an immutable digest (for example, \"image@sha256:<digest>\") so that a tag cannot be re-pointed at a different image."},
				},
			},
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
