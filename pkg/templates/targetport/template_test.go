package targetport

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/targetport/internal/params"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestTargetPort(t *testing.T) {
	suite.Run(t, new(TargetPortTestSuite))
}

type TargetPortTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *TargetPortTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *TargetPortTestSuite) AddServiceWithTargetPort(name string, targetPort intstr.IntOrString) {
	s.ctx.AddMockService(s.T(), name)
	s.ctx.ModifyService(s.T(), name, func(service *coreV1.Service) {
		service.Spec.Ports = append(service.Spec.Ports, coreV1.ServicePort{TargetPort: targetPort})
	})
}

func (s *TargetPortTestSuite) AddDeploymentWithPortName(name, portName string) {
	s.ctx.AddMockDeployment(s.T(), name)
	s.ctx.ModifyDeployment(s.T(), name, func(deployment *appsV1.Deployment) {
		deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, coreV1.Container{
			Name: name,
			Ports: []coreV1.ContainerPort{
				{Name: portName},
			},
		})
	})
}

func (s *TargetPortTestSuite) TestNoPorts() {
	s.ctx.AddMockService(s.T(), "test-empty")
	s.ctx.AddMockDeployment(s.T(), "test-empty")

	expectedDiagnostics := map[string][]diagnostic.Diagnostic{}
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param:                    params.Params{},
			Diagnostics:              expectedDiagnostics,
			ExpectInstantiationError: false,
		},
	})
}

func (s *TargetPortTestSuite) TestServices() {
	type testCase struct {
		TargetPort          intstr.IntOrString
		ExpectedDiagnostics []diagnostic.Diagnostic
	}

	testServices := map[string]testCase{
		"ok-number": {
			TargetPort:          intstr.FromInt(8080),
			ExpectedDiagnostics: nil,
		},
		"negative-number": {
			TargetPort: intstr.FromInt(-1),
			ExpectedDiagnostics: []diagnostic.Diagnostic{
				{Message: "port targetPort \"-1\" in service \"negative-number\" must be between 1 and 65535, inclusive"},
			},
		},
		"big-number": {
			TargetPort: intstr.FromInt(65536),
			ExpectedDiagnostics: []diagnostic.Diagnostic{
				{Message: "port targetPort \"65536\" in service \"big-number\" must be between 1 and 65535, inclusive"},
			},
		},
		"ok-string": {
			TargetPort:          intstr.FromString("port-name"),
			ExpectedDiagnostics: nil,
		},
		"ok-one-char-string": {
			TargetPort:          intstr.FromString("p"),
			ExpectedDiagnostics: nil,
		},
		"long-string": {
			TargetPort: intstr.FromString("long-name-over-15"),
			ExpectedDiagnostics: []diagnostic.Diagnostic{
				{Message: "port targetPort \"long-name-over-15\" in service \"long-string\" must be no more than 15 characters"},
			},
		},
		"empty-string": {
			TargetPort: intstr.FromString(""),
			ExpectedDiagnostics: []diagnostic.Diagnostic{
				{Message: "port targetPort \"\" in service \"empty-string\" must contain only alpha-numeric characters (a-z, 0-9), and hyphens (-)"},
				{Message: "port targetPort \"\" in service \"empty-string\" must contain at least one letter (a-z)"},
			},
		},
		"capital-string": {
			TargetPort: intstr.FromString("PORT-NAME"),
			ExpectedDiagnostics: []diagnostic.Diagnostic{
				{Message: "port targetPort \"PORT-NAME\" in service \"capital-string\" must contain only alpha-numeric characters (a-z, 0-9), and hyphens (-)"},
				{Message: "port targetPort \"PORT-NAME\" in service \"capital-string\" must contain at least one letter (a-z)"},
			},
		},
		"bad-char-string": {
			TargetPort: intstr.FromString("port_name"),
			ExpectedDiagnostics: []diagnostic.Diagnostic{
				{Message: "port targetPort \"port_name\" in service \"bad-char-string\" must contain only alpha-numeric characters (a-z, 0-9), and hyphens (-)"},
			},
		},
		"num-string": {
			TargetPort: intstr.FromString("8080"),
			ExpectedDiagnostics: []diagnostic.Diagnostic{
				{Message: "port targetPort \"8080\" in service \"num-string\" must contain at least one letter (a-z)"},
			},
		},
		"hyphen-start-string": {
			TargetPort: intstr.FromString("-port-name"),
			ExpectedDiagnostics: []diagnostic.Diagnostic{
				{Message: "port targetPort \"-port-name\" in service \"hyphen-start-string\" must not begin or end with a hyphen"},
			},
		},
		"hyphen-end-string": {
			TargetPort: intstr.FromString("port-name-"),
			ExpectedDiagnostics: []diagnostic.Diagnostic{
				{Message: "port targetPort \"port-name-\" in service \"hyphen-end-string\" must not begin or end with a hyphen"},
			},
		},
		"double-hyphen-string": {
			TargetPort: intstr.FromString("port--name"),
			ExpectedDiagnostics: []diagnostic.Diagnostic{
				{Message: "port targetPort \"port--name\" in service \"double-hyphen-string\" must not contain consecutive hyphens"},
			},
		},
	}

	expectedDiagnostics := map[string][]diagnostic.Diagnostic{}
	for serviceName, info := range testServices {
		s.AddServiceWithTargetPort(serviceName, info.TargetPort)
		expectedDiagnostics[serviceName] = append(expectedDiagnostics[serviceName], info.ExpectedDiagnostics...)
	}

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param:                    params.Params{},
			Diagnostics:              expectedDiagnostics,
			ExpectInstantiationError: false,
		},
	})
}

func (s *TargetPortTestSuite) TestDeployments() {
	type testCase struct {
		TargetPort          string
		ExpectedDiagnostics []diagnostic.Diagnostic
	}

	testDeployments := map[string]testCase{
		"ok-string": {
			TargetPort:          "port-name",
			ExpectedDiagnostics: nil,
		},
		"ok-one-char-string": {
			TargetPort:          "p",
			ExpectedDiagnostics: nil,
		},
		"long-string": {
			TargetPort: "long-name-over-15",
			ExpectedDiagnostics: []diagnostic.Diagnostic{
				{Message: "port name \"long-name-over-15\" in container \"long-string\" must be no more than 15 characters"},
			},
		},
		"capital-string": {
			TargetPort: "PORT-NAME",
			ExpectedDiagnostics: []diagnostic.Diagnostic{
				{Message: "port name \"PORT-NAME\" in container \"capital-string\" must contain only alpha-numeric characters (a-z, 0-9), and hyphens (-)"},
				{Message: "port name \"PORT-NAME\" in container \"capital-string\" must contain at least one letter (a-z)"},
			},
		},
		"bad-char-string": {
			TargetPort: "port_name",
			ExpectedDiagnostics: []diagnostic.Diagnostic{
				{Message: "port name \"port_name\" in container \"bad-char-string\" must contain only alpha-numeric characters (a-z, 0-9), and hyphens (-)"},
			},
		},
		"num-string": {
			TargetPort: "8080",
			ExpectedDiagnostics: []diagnostic.Diagnostic{
				{Message: "port name \"8080\" in container \"num-string\" must contain at least one letter (a-z)"},
			},
		},
		"hyphen-start-string": {
			TargetPort: "-port-name",
			ExpectedDiagnostics: []diagnostic.Diagnostic{
				{Message: "port name \"-port-name\" in container \"hyphen-start-string\" must not begin or end with a hyphen"},
			},
		},
		"hyphen-end-string": {
			TargetPort: "port-name-",
			ExpectedDiagnostics: []diagnostic.Diagnostic{
				{Message: "port name \"port-name-\" in container \"hyphen-end-string\" must not begin or end with a hyphen"},
			},
		},
		"double-hyphen-string": {
			TargetPort: "port--name",
			ExpectedDiagnostics: []diagnostic.Diagnostic{
				{Message: "port name \"port--name\" in container \"double-hyphen-string\" must not contain consecutive hyphens"},
			},
		},
	}

	expectedDiagnostics := map[string][]diagnostic.Diagnostic{}
	for deploymentName, info := range testDeployments {
		s.AddDeploymentWithPortName(deploymentName, info.TargetPort)
		expectedDiagnostics[deploymentName] = append(expectedDiagnostics[deploymentName], info.ExpectedDiagnostics...)
	}

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param:                    params.Params{},
			Diagnostics:              expectedDiagnostics,
			ExpectInstantiationError: false,
		},
	})
}
