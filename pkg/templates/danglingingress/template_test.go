package danglingingress

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/danglingingress/internal/params"
	coreV1 "k8s.io/api/core/v1"
	networkingV1 "k8s.io/api/networking/v1"
)

const (
	ingressName1 = "ingress-1"
	ingressName2 = "ingress-2"
	serviceName1 = "service-1"
	serviceName2 = "service-2"
)

const (
	port int32 = 80
)

func TestDanglingIngress(t *testing.T) {
	suite.Run(t, new(DanglingIngressTestSuite))
}

type DanglingIngressTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *DanglingIngressTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *DanglingIngressTestSuite) addIngress(name string, serviceSelectors []string) {
	s.ctx.AddMockIngress(s.T(), name)

	var paths []networkingV1.HTTPIngressPath
	for _, s := range serviceSelectors {
		paths = append(paths, networkingV1.HTTPIngressPath{
			Backend: networkingV1.IngressBackend{
				Service: &networkingV1.IngressServiceBackend{
					Name: s,
					Port: networkingV1.ServiceBackendPort{
						Number: port,
					},
				},
			},
		})
	}

	s.ctx.ModifyIngress(s.T(), name, func(ingress *networkingV1.Ingress) {
		ingress.Spec.Rules = []networkingV1.IngressRule{{
			IngressRuleValue: networkingV1.IngressRuleValue{
				HTTP: &networkingV1.HTTPIngressRuleValue{
					Paths: paths,
				},
			},
		}}
	})
}

func (s *DanglingIngressTestSuite) addIngressWithNamedPorts(name, serviceSelector, portName string) {
	s.ctx.AddMockIngress(s.T(), name)

	paths := []networkingV1.HTTPIngressPath{{
		Backend: networkingV1.IngressBackend{
			Service: &networkingV1.IngressServiceBackend{
				Name: serviceSelector,
				Port: networkingV1.ServiceBackendPort{
					Name: portName,
				},
			},
		},
	}}

	s.ctx.ModifyIngress(s.T(), name, func(ingress *networkingV1.Ingress) {
		ingress.Spec.Rules = []networkingV1.IngressRule{{
			IngressRuleValue: networkingV1.IngressRuleValue{
				HTTP: &networkingV1.HTTPIngressRuleValue{
					Paths: paths,
				},
			},
		}}
	})
}

func (s *DanglingIngressTestSuite) addIngressWithDefaultBackend(name, serviceName string, serviceSelectors []string) {
	s.addIngress(name, serviceSelectors)

	s.ctx.ModifyIngress(s.T(), name, func(ingress *networkingV1.Ingress) {
		ingress.Spec.DefaultBackend = &networkingV1.IngressBackend{
			Service: &networkingV1.IngressServiceBackend{
				Name: serviceName,
				Port: networkingV1.ServiceBackendPort{
					Number: port,
				},
			},
		}
	})
}

func (s *DanglingIngressTestSuite) addService(name string, port int32) {
	s.addServiceWithPortName(name, port, "")
}

func (s *DanglingIngressTestSuite) addServiceWithPortName(name string, port int32, portName string) {
	s.ctx.AddMockService(s.T(), name)
	s.ctx.ModifyService(s.T(), name, func(s *coreV1.Service) {
		s.Spec.Ports = append(s.Spec.Ports, coreV1.ServicePort{
			Name: portName,
			Port: port,
		})
	})
}

func (s *DanglingIngressTestSuite) TestIngressWithOutServicesPasses() {
	s.ctx.AddMockIngress(s.T(), ingressName1)
	s.ctx.AddMockDeployment(s.T(), "deployment")

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				ingressName1: {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingIngressTestSuite) TestIngressFailsWithNoService() {
	s.addIngress(ingressName1, []string{serviceName1})
	s.addService("something-other-service", port)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				ingressName1: {{
					Message: "no service found matching ingress label (service-1), port 80",
				}},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingIngressTestSuite) TestIngressWithServicePasses() {
	s.addIngress(ingressName1, []string{serviceName1})
	s.addService(serviceName1, port)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				ingressName1: {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingIngressTestSuite) TestIngressWithServiceMissmatchPortfails() {
	s.addIngress(ingressName1, []string{serviceName1})
	s.addIngressWithNamedPorts(ingressName2, serviceName1, "port1")
	s.addServiceWithPortName(serviceName1, 9000, "port2")

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				ingressName1: {{Message: "no service found matching ingress label (service-1), port 80"}},
				ingressName2: {{Message: "no service found matching ingress label (service-1), port port1"}},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingIngressTestSuite) TestIngresssWithMissmatchPortNumbersWillFail() {
	s.addIngress(ingressName1, []string{serviceName1})
	s.addService(serviceName1, 9000)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				ingressName1: {{Message: "no service found matching ingress label (service-1), port 80"}},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingIngressTestSuite) TestIngressWithMatchingPortNameWillPass() {
	const portName = "port1"
	s.addIngressWithNamedPorts(ingressName1, serviceName1, portName)
	s.addServiceWithPortName(serviceName1, port, portName)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				ingressName1: {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingIngressTestSuite) TestIngressWithAnyMissingServiceFails() {
	s.addIngress(ingressName1, []string{serviceName1, "not-existent-service"})
	s.addService(serviceName1, port)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				ingressName1: {{Message: "no service found matching ingress label (not-existent-service), port 80"}},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingIngressTestSuite) TestIngressWillReportAllMissingServices() {
	s.addIngress(ingressName1, []string{serviceName1, "not-existent-service"})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				ingressName1: {
					{Message: "no service found matching ingress label (service-1), port 80"},
					{Message: "no service found matching ingress label (not-existent-service), port 80"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingIngressTestSuite) TestIngressWillFindAllServicesIfTheyExist() {
	serviceName2 := "service-2"
	s.addIngress(ingressName1, []string{serviceName1, serviceName2})
	s.addService(serviceName1, port)
	s.addService(serviceName2, port)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				ingressName1: {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingIngressTestSuite) TestDanglingIngressWillPassWithMultipleIngresses() {
	s.addIngress(ingressName1, []string{serviceName1})
	s.addIngress(ingressName2, []string{serviceName2})
	s.addService(serviceName1, port)
	s.addService(serviceName2, port)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				ingressName1: {},
				ingressName2: {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingIngressTestSuite) TestDanglingIngressWillFailWithOneDangling() {
	s.addIngress(ingressName1, []string{serviceName1})
	s.addIngress(ingressName2, []string{serviceName2})
	s.addService(serviceName1, port)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				ingressName1: {},
				ingressName2: {
					{Message: "no service found matching ingress label (service-2), port 80"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingIngressTestSuite) TestIngressWithDefaultBackendServiceExistsPasses() {
	s.addIngressWithDefaultBackend(ingressName1, serviceName1, []string{})
	s.addService(serviceName1, port)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				ingressName1: {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingIngressTestSuite) TestIngressWithDefaultBackendAndRulesServiceExistsPasses() {
	s.addIngressWithDefaultBackend(ingressName1, serviceName1, []string{serviceName1})
	s.addService(serviceName1, port)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				ingressName1: {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingIngressTestSuite) TestIngressWithDefaultBackendServiceMissingFails() {
	s.addIngressWithDefaultBackend(ingressName1, serviceName1, []string{})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				ingressName1: {{Message: "no service found matching ingress label (service-1), port 80"}},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingIngressTestSuite) TestIngressWithDefaultBackendAndRulesServiceMissingFails() {
	s.addIngressWithDefaultBackend(ingressName2, serviceName2, []string{serviceName1})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				ingressName2: {
					{Message: "no service found matching ingress label (service-2), port 80"},
					{Message: "no service found matching ingress label (service-1), port 80"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingIngressTestSuite) TestIngressWithResoucesInsteadofServicesPasses() {
	s.ctx.AddMockIngress(s.T(), ingressName1)
	resouceName := "resource"
	apiGroup := "v2"

	s.ctx.ModifyIngress(s.T(), ingressName1, func(ingress *networkingV1.Ingress) {
		ingress.Spec.Rules = []networkingV1.IngressRule{{
			IngressRuleValue: networkingV1.IngressRuleValue{
				HTTP: &networkingV1.HTTPIngressRuleValue{
					Paths: []networkingV1.HTTPIngressPath{
						{
							Backend: networkingV1.IngressBackend{
								Resource: &coreV1.TypedLocalObjectReference{
									APIGroup: &apiGroup,
									Kind:     "Pod",
									Name:     resouceName,
								},
							},
						},
					},
				},
			},
		}}
	})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				ingressName1: {},
			},
			ExpectInstantiationError: false,
		},
	})
}
