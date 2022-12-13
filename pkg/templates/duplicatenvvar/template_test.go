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

func (s *DuplicateEnvVarTestSuite) TestDeploymentWith() {
	const targetName = "deployment01"
	testCases := []struct {
		name      string
		container v1.Container
		expected  map[string][]diagnostic.Diagnostic
	}{
		{
			name: "WithNoDuplicatesShouldPass",
			container: v1.Container{
				Env: []v1.EnvVar{
					{Name: "ENV_1"},
					{Name: "ENV_2"},
				},
			},
			expected: nil,
		},
		{
			name: "DuplicatesShouldReportErrors",
			container: v1.Container{
				Env: []v1.EnvVar{
					{Name: "ENV_1"},
					{Name: "ENV_1"},
				},
			},
			expected: map[string][]diagnostic.Diagnostic{
				targetName: {
					// Ensure we only get one error for ENV_1
					{Message: "Duplicate environment variable ENV_1 in container \"\" found"},
				},
			},
		},
		{
			name: "MultipleDuplicatesShouldReportThemAllButOnlyOncePerEnv",
			container: v1.Container{
				Name: "container",
				Env: []v1.EnvVar{
					{Name: "ENV_1"},
					{Name: "ENV_1"},
					{Name: "ENV_1"},
					{Name: "ENV_2"},
					{Name: "ENV_2"},
					{Name: "ENV_3"},
				},
			},
			expected: map[string][]diagnostic.Diagnostic{
				targetName: {
					// Ensure we only get one error for ENV_1
					{Message: "Duplicate environment variable ENV_1 in container \"container\" found"},
					{Message: "Duplicate environment variable ENV_2 in container \"container\" found"},
				},
			},
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.ctx.AddMockDeployment(s.T(), targetName)
			s.ctx.AddContainerToDeployment(s.T(), targetName, tc.container)
			s.Validate(s.ctx, []templates.TestCase{{
				Param:       params.Params{},
				Diagnostics: tc.expected,
			}})
		})
	}
}
