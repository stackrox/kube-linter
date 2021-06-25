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

func (s *ContainerImageTestSuite) TestImproperContainerTag() {
	const (
		depWithLatestAsContainerImageTag = "dep-with-latest-as-container-image-tag"
	)

	s.addDeploymentWithContainerImage(depWithLatestAsContainerImageTag, "example.com/test:latest")

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				BlockList: []string{".*:(latest)$"},
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				depWithLatestAsContainerImageTag: {
					{Message: "The container \"test-container\" is using a floating image tag, \"example.com/test:latest\"."},
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
	})
}
