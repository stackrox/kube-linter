package imagepullpolicy

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/imagepullpolicy/internal/params"

	v1 "k8s.io/api/core/v1"
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

func (s *ContainerImageTestSuite) addDeploymentWithContainerImage(name string, pullPolicy v1.PullPolicy) {
	s.ctx.AddMockDeployment(s.T(), name)
	s.ctx.AddContainerToDeployment(s.T(), name, v1.Container{Name: "test-container", ImagePullPolicy: pullPolicy})
}

func (s *ContainerImageTestSuite) TestImaPolicy() {
	const (
		alwaysDep       = "deployment-with-always-pull-policy"
		ifNotPresentDep = "deployment-with-if-not-present-pull-policy"
		neverDep        = "deployment-with-never-pull-policy"
	)

	s.addDeploymentWithContainerImage(alwaysDep, v1.PullAlways)
	s.addDeploymentWithContainerImage(ifNotPresentDep, v1.PullIfNotPresent)
	s.addDeploymentWithContainerImage(neverDep, v1.PullNever)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				ForbiddenPolicies: []string{"Always"},
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				alwaysDep: {
					{Message: "container \"test-container\" has imagePullPolicy set to Always"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}
