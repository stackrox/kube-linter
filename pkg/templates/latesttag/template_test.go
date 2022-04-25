package latesttag

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/latesttag/internal/params"

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
		depWithLatestAsContainerImageTag = "dep-with-latest-as-container-image-tag"
		depWithNotAllowedImageRegistry   = "dep-with-not-allowed-image-registry"
	)

	s.addDeploymentWithContainerImage(depWithLatestAsContainerImageTag, "example.com/test:latest")
	s.addDeploymentWithContainerImage(depWithNotAllowedImageRegistry, "test.com/test:v1.0.0")

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				BlockList: []string{".*:(latest)$"},
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				depWithLatestAsContainerImageTag: {
					{Message: "The container \"test-container\" is using an invalid container image, \"example.com/test:latest\". Please use images that are not blocked by the `BlockList` criteria : [\".*:(latest)$\"]"},
				},
			},
			ExpectInstantiationError: false,
		},
		{
			Param: params.Params{
				AllowList: []string{"^(example.com)"},
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				depWithNotAllowedImageRegistry: {
					{Message: "The container \"test-container\" is using an invalid container image, \"test.com/test:v1.0.0\". Please use images that satisfies the `AllowList` criteria : [\"^(example.com)\"]"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *ContainerImageTestSuite) TestAcceptableContainerImage() {
	const (
		depWithLatestAsContainerImageName = "dep-with-latest-as-container-image-name"
		depWithAcceptableContainerImage   = "dep-with-acceptable-container-image"
	)

	s.addDeploymentWithContainerImage(depWithLatestAsContainerImageName, "example.com/latest:v1.0.0")
	s.addDeploymentWithContainerImage(depWithAcceptableContainerImage, "example.com/test:v1.0.0")

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				BlockList: []string{".*:(latest)$"},
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				depWithLatestAsContainerImageName: nil,
				depWithAcceptableContainerImage:   nil,
			},
			ExpectInstantiationError: false,
		},
		{
			Param: params.Params{
				AllowList: []string{"^(example.com)"},
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				depWithAcceptableContainerImage: nil,
			},
			ExpectInstantiationError: false,
		},
	})
}
