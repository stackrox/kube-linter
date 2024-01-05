package readinessport

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/readinessport/internal/params"
	v1 "k8s.io/api/core/v1"
)

func TestMissingReadinessPort(t *testing.T) {
	suite.Run(t, new(MissingReadinessPort))
}

type MissingReadinessPort struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *MissingReadinessPort) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *MissingReadinessPort) TestDeploymentWith() {
	const targetName = "deployment01"
	testCases := []struct {
		name      string
		container v1.Container
		expected  map[string][]diagnostic.Diagnostic
	}{
		{
			name:      "NoReadinessProbe",
			container: v1.Container{},
			expected:  nil,
		},
		{
			name: "NoReadinessProbeExecIgnored",
			container: v1.Container{
				ReadinessProbe: &v1.Probe{
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
				ReadinessProbe: &v1.Probe{
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
				ReadinessProbe: &v1.Probe{
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
