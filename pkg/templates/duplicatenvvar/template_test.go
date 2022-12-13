package duplicateenvvar

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/duplicatenvvar/internal/params"
	v1 "k8s.io/api/core/v1"
)

func TestDuplicateEnvVar(t *testing.T) {
	suite.Run(t, new(DuplicateEnvVarTestSuite))
}

type DuplicateEnvVarTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *DuplicateEnvVarTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *DuplicateEnvVarTestSuite) TestDeploymentWithNoDuplicatesPass() {
	const targetName = "deployment01"

	s.ctx.AddMockDeployment(s.T(), targetName)
	s.ctx.AddContainerToDeployment(s.T(), targetName, v1.Container{
		Env: []v1.EnvVar{
			{
				Name: "ENV_1",
			},
			{
				Name: "ENV_2",
			},
		},
	})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				targetName: nil,
			},
		},
	})
}

func (s *DuplicateEnvVarTestSuite) TestDeploymentWillReportErrorsWithDuplicates() {
	const targetName = "deployment02"

	s.ctx.AddMockDeployment(s.T(), targetName)
	s.ctx.AddContainerToDeployment(s.T(), targetName, v1.Container{
		Name: "containerName",
		Env: []v1.EnvVar{
			{
				Name: "ENV_1",
			},
			{
				Name: "ENV_1",
			},
		},
	})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				targetName: {
					{Message: "Duplicate environment variable ENV_1 in container \"containerName\" found"},
				},
			},
		},
	})
}

func (s *DuplicateEnvVarTestSuite) TestDeploymentWillReportAllDuplicates() {
	const targetName = "deployment02"

	s.ctx.AddMockDeployment(s.T(), targetName)
	s.ctx.AddContainerToDeployment(s.T(), targetName, v1.Container{
		Name: "container",
		Env: []v1.EnvVar{
			{Name: "ENV_1"},
			{Name: "ENV_1"},
			{Name: "ENV_1"},
			{Name: "ENV_2"},
			{Name: "ENV_2"},
		},
	})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				targetName: {
					// Ensure we only get one error for ENV_1
					{Message: "Duplicate environment variable ENV_1 in container \"container\" found"},
					{Message: "Duplicate environment variable ENV_2 in container \"container\" found"},
				},
			},
		},
	})
}
