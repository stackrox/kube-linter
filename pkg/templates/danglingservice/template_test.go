package danglingservice

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/danglingservice/internal/params"

	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

const (
	pod1        = "pod1"
	pod2        = "pod2"
	serviceNone = "service-matches-none"
	service1    = "service-matches-pod1"
	service2    = "service-matches-pod2"
)

var emptyLabels = map[string]string{} //empty selector

var labelselector1 = map[string]string{"app": "pod1-test"}

var labelselector2 = map[string]string{"app": "pod2-test"}

var labelSelectorMulti = map[string]string{"app": "pod2-test", "env": "staging"}

func TestDanglingService(t *testing.T) {
	suite.Run(t, new(DanglingServiceTestSuite))
}

type DanglingServiceTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *DanglingServiceTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *DanglingServiceTestSuite) AddService(name string, podLabels map[string]string) {
	s.ctx.AddMockService(s.T(), name)
	s.ctx.ModifyService(s.T(), name, func(service *v1.Service) {
		service.Spec.Selector = podLabels
	})
}

func (s *DanglingServiceTestSuite) AddDeploymentWithLabels(name string, labels map[string]string) {
	s.ctx.AddMockDeployment(s.T(), name)
	s.ctx.ModifyDeployment(s.T(), name, func(deployment *appsV1.Deployment) {
		deployment.Spec.Template.Labels = labels
	})
}

func (s *DanglingServiceTestSuite) TestServiceEmptySelectorMatchesNoPods() {
	s.AddDeploymentWithLabels(pod1, labelselector1)
	s.AddDeploymentWithLabels(pod2, labelselector2)
	s.AddService(serviceNone, emptyLabels)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				serviceNone: {{Message: "service has no selector specified"}},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingServiceTestSuite) TestNoDanglingServices() {
	s.AddDeploymentWithLabels(pod1, labelselector1)
	s.AddDeploymentWithLabels(pod2, labelselector2)
	s.AddService(service1, labelselector1)
	s.AddService(service2, labelselector2)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				service1: {},
				service2: {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingServiceTestSuite) TestOneServiceIsDangling() {
	s.AddDeploymentWithLabels(pod2, labelselector2)
	s.AddService(service1, labelselector1)
	s.AddService(service2, labelselector2)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				service1: {{Message: fmt.Sprintf("no pods found matching service labels (%v)", labelselector1)}},
				service2: {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingServiceTestSuite) TestMatchingWithIgnoredLabel() {
	s.AddDeploymentWithLabels(pod2, labelselector2) // only app label
	s.AddService(service1, labelselector1)          // app label but doesn't match
	s.AddService(service2, labelSelectorMulti)      // app and env labels

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{IgnoredLabels: []string{"env"}}, // Ignore missing env label
			Diagnostics: map[string][]diagnostic.Diagnostic{
				service1: {{Message: fmt.Sprintf("no pods found matching service labels (%v)", labelselector1)}},
				service2: {},
			},
			ExpectInstantiationError: false,
		},
	})
}
