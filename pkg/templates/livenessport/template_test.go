package livenessport

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/livenessport/internal/params"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestMissingLivenessPort(t *testing.T) {
	suite.Run(t, new(MissingLivenessPort))
}

type MissingLivenessPort struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *MissingLivenessPort) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *MissingLivenessPort) TestDeploymentWith() {
	const targetName = "deployment01"
	testCases := []struct {
		name      string
		container v1.Container
		expected  map[string][]diagnostic.Diagnostic
	}{
		{
			name:      "NoLivenessProbe",
			container: v1.Container{},
			expected:  nil,
		},
		{
			name: "NoLivenessProbeExecIgnored",
			container: v1.Container{
				LivenessProbe: &v1.Probe{
					ProbeHandler: v1.ProbeHandler{
						Exec: &v1.ExecAction{},
					},
				},
			},
			expected: nil,
		},
		{
			name: "MatchinPortInt",
			container: v1.Container{
				Ports: []v1.ContainerPort{
					{
						Name:          "http",
						ContainerPort: 8080,
						Protocol:      v1.ProtocolTCP,
					},
				},
				LivenessProbe: &v1.Probe{
					ProbeHandler: v1.ProbeHandler{
						HTTPGet: &v1.HTTPGetAction{
							Port: intstr.FromInt(8080),
						},
					},
				},
			},
			expected: nil,
		},
		{
			name: "MatchinPortStr",
			container: v1.Container{
				Ports: []v1.ContainerPort{
					{
						Name:          "http",
						ContainerPort: 8080,
					},
				},
				LivenessProbe: &v1.Probe{
					ProbeHandler: v1.ProbeHandler{
						HTTPGet: &v1.HTTPGetAction{
							Port: intstr.FromString("http"),
						},
					},
				},
			},
			expected: nil,
		},
		{
			name: "MatchinPortStrTCPSSocket",
			container: v1.Container{
				Ports: []v1.ContainerPort{
					{
						Name:          "http",
						ContainerPort: 8080,
					},
				},
				LivenessProbe: &v1.Probe{
					ProbeHandler: v1.ProbeHandler{
						TCPSocket: &v1.TCPSocketAction{
							Port: intstr.FromString("http"),
						},
					},
				},
			},
			expected: nil,
		},
		{
			name: "MismaptchPort",
			container: v1.Container{
				Name: "container",
				Ports: []v1.ContainerPort{
					{
						Name:          "http",
						ContainerPort: 8080,
					},
				},
				LivenessProbe: &v1.Probe{
					ProbeHandler: v1.ProbeHandler{
						HTTPGet: &v1.HTTPGetAction{
							Port: intstr.FromInt(9999),
						},
					},
				},
			},
			expected: map[string][]diagnostic.Diagnostic{
				targetName: {
					{Message: "container \"container\" does not expose port 9999 for the HTTPGet"},
				},
			},
		},
		{
			name: "MismaptchPort",
			container: v1.Container{
				Name: "container",
				Ports: []v1.ContainerPort{
					{
						Name:          "http",
						ContainerPort: 8080,
					},
				},
				LivenessProbe: &v1.Probe{
					ProbeHandler: v1.ProbeHandler{
						HTTPGet: &v1.HTTPGetAction{
							Port: intstr.FromString("healthcheck"),
						},
					},
				},
			},
			expected: map[string][]diagnostic.Diagnostic{
				targetName: {
					{Message: "container \"container\" does not expose port healthcheck for the HTTPGet"},
				},
			},
		},
		{
			name: "MatchinPortUDP",
			container: v1.Container{
				Name: "container",
				Ports: []v1.ContainerPort{
					{
						Name:          "http",
						ContainerPort: 8080,
						Protocol:      v1.ProtocolUDP,
					},
				},
				LivenessProbe: &v1.Probe{
					ProbeHandler: v1.ProbeHandler{
						HTTPGet: &v1.HTTPGetAction{
							Port: intstr.FromInt(8080),
						},
					},
				},
			},
			expected: map[string][]diagnostic.Diagnostic{
				targetName: {
					{Message: "container \"container\" does not expose port 8080 for the HTTPGet"},
				},
			},
		},
		{
			name: "MismaptchPortTCPSocket",
			container: v1.Container{
				Name: "container",
				Ports: []v1.ContainerPort{
					{
						Name:          "http",
						ContainerPort: 8080,
					},
				},
				LivenessProbe: &v1.Probe{
					ProbeHandler: v1.ProbeHandler{
						TCPSocket: &v1.TCPSocketAction{
							Port: intstr.FromString("socket"),
						},
					},
				},
			},
			expected: map[string][]diagnostic.Diagnostic{
				targetName: {
					{Message: "container \"container\" does not expose port socket for the TCPSocket"},
				},
			},
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
				LivenessProbe: &v1.Probe{
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
				LivenessProbe: &v1.Probe{
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
