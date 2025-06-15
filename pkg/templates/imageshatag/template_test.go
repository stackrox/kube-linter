package imageshatag

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/imageshatag/internal/params"

	v1 "k8s.io/api/core/v1"
)

var (
	containerName = "test-container"
)

func TestContainerImage(t *testing.T) {
	suite.Run(t, new(ContainerImageTestSuite))
}

type ContainerImageTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *ContainerImageTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *ContainerImageTestSuite) addDeploymentWithContainerImage(name, containerImage string) {
	s.ctx.AddMockDeployment(s.T(), name)
	s.ctx.AddContainerToDeployment(s.T(), name, v1.Container{Name: containerName, Image: containerImage})
}

func (s *ContainerImageTestSuite) TestImproperContainerImage() {
	const (
		depwithNotAllowedImageTag = "dep-with-not-allowed-image-tag"
	)

	s.addDeploymentWithContainerImage(depwithNotAllowedImageTag, "example.com/test:latest")

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				AllowList: []string{".*:[a-fA-F0-9]{64}$"},
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				depwithNotAllowedImageTag: {
					{Message: "The container \"test-container\" is using an invalid container image, \"example.com/test:latest\". Please reference the image using a SHA256 tag."},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *ContainerImageTestSuite) TestAcceptableContainerImage() {
	const (
		depWithAcceptableImageTag = "dep-with-acceptable-container-image"
	)

	s.addDeploymentWithContainerImage(depWithAcceptableImageTag, "example.com/latest@sha256:75bf9b911b6481dcf29f7942240d1555adaa607eec7fc61bedb7f624f87c36d4")
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				AllowList: []string{".*:[a-fA-F0-9]{64}$"},
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				depWithAcceptableImageTag: nil,
			},
			ExpectInstantiationError: false,
		},
	})
}
