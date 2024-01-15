package startupport

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/startupport/internal/params"
	v1 "k8s.io/api/core/v1"
)

func TestMissingStartUpPort(t *testing.T) {
	suite.Run(t, new(MissingStartUpPort))
}

type MissingStartUpPort struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *MissingStartUpPort) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *MissingStartUpPort) TestDeploymentWith() {
	const targetName = "deployment01"
	testCases := []struct {
		name      string
		container v1.Container
		expected  map[string][]diagnostic.Diagnostic
	}{
		{
			name:      "NoStartUpProbe",
			container: v1.Container{},
			expected:  nil,
		},
		{
			name: "NoStartUpProbeExecIgnored",
			container: v1.Container{
				StartupProbe: &v1.Probe{
					ProbeHandler: v1.ProbeHandler{
						Exec: &v1.ExecAction{},
					},
				},
			},
			expected: nil,
		},
		{
			name: "GrpcCheckWillPass",
			container: v1.Container{
				Name: "container",
				Ports: []v1.ContainerPort{
					{
						Name:          "http",
						ContainerPort: 8080,
					},
				},
				StartupProbe: &v1.Probe{
					ProbeHandler: v1.ProbeHandler{
						GRPC: &v1.GRPCAction{
							Port: 8080,
						},
					},
				},
			},
		},
		{
			name: "GrpcPortMissmatch",
			container: v1.Container{
				Name: "container",
				Ports: []v1.ContainerPort{
					{
						Name:          "http",
						ContainerPort: 8080,
					},
				},
				StartupProbe: &v1.Probe{
					ProbeHandler: v1.ProbeHandler{
						GRPC: &v1.GRPCAction{
							Port: 9999,
						},
					},
				},
			},
			expected: map[string][]diagnostic.Diagnostic{
				targetName: {
					{Message: "container \"container\" does not expose port 9999 for the GRPC check"},
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
