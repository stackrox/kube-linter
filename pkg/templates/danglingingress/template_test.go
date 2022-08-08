package danglingingress

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/danglingingress/internal/params"
	networkingV1 "k8s.io/api/networking/v1"
)

const (
	ingressName1 = "ingress-1"
	ingressName2 = "ingress-2"
	serviceName1 = "service-1"
	serviceName2 = "service-2"
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
				},
			},
		})
	}

	s.ctx.ModifyIngess(s.T(), name, func(ingress *networkingV1.Ingress) {
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

	s.ctx.ModifyIngess(s.T(), name, func(ingress *networkingV1.Ingress) {
		ingress.Spec.DefaultBackend = &networkingV1.IngressBackend{
			Service: &networkingV1.IngressServiceBackend{
				Name: serviceName,
			},
		}
	})
}

func (s *DanglingIngressTestSuite) addService(name string) {
	s.ctx.AddMockService(s.T(), name)
}

func (s *DanglingIngressTestSuite) TestIngressWithBackendFails() {
	s.ctx.AddMockIngress(s.T(), ingressName1)
	s.ctx.AddMockDeployment(s.T(), "deployment")

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				ingressName1: {{
					Message: "ingress has no backend specified",
				}},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingIngressTestSuite) TestIngressFailsWithNoService() {
	s.addIngress(ingressName1, []string{serviceName1})
	s.addService("something-other-service")

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				ingressName1: {{
					Message: "no service found matching ingress labels (service-1)",
				}},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingIngressTestSuite) TestIngressWithServicePasses() {
	s.addIngress(ingressName1, []string{serviceName1})
	s.addService(serviceName1)

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
	s.addService(serviceName1)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				ingressName1: {{Message: "no service found matching ingress labels (not-existent-service)"}},
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
					{Message: "no service found matching ingress labels (service-1)"},
					{Message: "no service found matching ingress labels (not-existent-service)"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingIngressTestSuite) TestIngressWillFindAllServicesIfTheyExist() {
	serviceName2 := "service-2"
	s.addIngress(ingressName1, []string{serviceName1, serviceName2})
	s.addService(serviceName1)
	s.addService(serviceName2)

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
	s.addService(serviceName1)
	s.addService(serviceName2)

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
	s.addService(serviceName1)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				ingressName1: {},
				ingressName2: {
					{Message: "no service found matching ingress labels (service-2)"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingIngressTestSuite) TestIngressWithDefaultBackendServiceExistsPasses() {
	s.addIngressWithDefaultBackend(ingressName1, serviceName1, []string{})
	s.addService(serviceName1)

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
	s.addService(serviceName1)

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
				ingressName1: {{Message: "no service found matching ingress labels (service-1)"}},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingIngressTestSuite) TestIngressWithDefaultBackendAndRulesServiceMissingFails() {
	s.addIngressWithDefaultBackend(ingressName1, serviceName1, []string{serviceName2})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				ingressName1: {
					{Message: "no service found matching ingress labels (service-1)"},
					{Message: "no service found matching ingress labels (service-2)"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}
