package danglingservicemonitor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/danglingservicemonitor/internal/params"

	k8sMonitoring "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	service1           = "service1"
	service2           = "service2"
	servicemonitor1    = "servicemonitor-matches-service1"
	servicemonitor2    = "servicemonitor-matches-service2"
	servicemonitorNone = "servicemonitor-matches-none"
)

var emptyLabelSelector = metaV1.LabelSelector{
	MatchLabels: map[string]string{},
}

var labelselector1 = metaV1.LabelSelector{
	MatchLabels: map[string]string{"app": "service1-test"},
}

var labelselector2 = metaV1.LabelSelector{
	MatchLabels: map[string]string{"app": "service2-test"},
}

var invalidlabelselector = metaV1.LabelSelector{
	MatchLabels: map[string]string{"-incorrect-labelselector-and-it-is-too-long-for-kubernetes-today-": "test"},
}
var namespaceselector1 = k8sMonitoring.NamespaceSelector{
	MatchNames: []string{"test1"},
}

var namespace1 = "test1"

var namespace2 = "test2"

func TestDanglingServiceMonitor(t *testing.T) {
	suite.Run(t, new(DanglingServiceMonitorTestSuite))
}

type DanglingServiceMonitorTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *DanglingServiceMonitorTestSuite) SetupTest() {
	s.Init("dangling-servicemonitor")
	s.ctx = mocks.NewMockContext()
}

func (s *DanglingServiceMonitorTestSuite) AddServiceMonitorWithLabelSelector(name string, labelSelector metaV1.LabelSelector) {
	s.ctx.AddMockServiceMonitor(s.T(), name)
	s.ctx.ModifyServiceMonitor(s.T(), name, func(servicemonitor *k8sMonitoring.ServiceMonitor) {
		servicemonitor.Spec.Selector = labelSelector
	})
}

func (s *DanglingServiceMonitorTestSuite) AddServiceMonitorWithNamespaceSelector(name string, namespaceSelector k8sMonitoring.NamespaceSelector) {
	s.ctx.AddMockServiceMonitor(s.T(), name)
	s.ctx.ModifyServiceMonitor(s.T(), name, func(servicemonitor *k8sMonitoring.ServiceMonitor) {
		servicemonitor.Spec.NamespaceSelector = namespaceSelector
	})
}

func (s *DanglingServiceMonitorTestSuite) AddServiceWithLabels(name string, labels *metaV1.LabelSelector) {
	s.ctx.AddMockService(s.T(), name)
	s.ctx.ModifyService(s.T(), name, func(service *coreV1.Service) {
		service.Labels = labels.MatchLabels
	})
}

func (s *DanglingServiceMonitorTestSuite) AddServiceWithNamespace(name, namespace string) {
	s.ctx.AddMockService(s.T(), name)
	s.ctx.ModifyService(s.T(), name, func(service *coreV1.Service) {
		service.Namespace = namespace
	})
}

func (s *DanglingServiceMonitorTestSuite) TestServiceMonitorEmpty() {
	s.AddServiceWithLabels(service1, &labelselector1)
	s.AddServiceWithLabels(service2, &labelselector2)
	s.AddServiceMonitorWithLabelSelector(servicemonitorNone, emptyLabelSelector)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				servicemonitorNone: {{Message: "service monitor has no selector specified"}},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingServiceMonitorTestSuite) TestNoDanglingServiceMonitors() {
	s.AddServiceWithLabels(service1, &labelselector1)
	s.AddServiceWithLabels(service2, &labelselector2)
	s.AddServiceMonitorWithLabelSelector(servicemonitor1, labelselector1)
	s.AddServiceMonitorWithLabelSelector(servicemonitor2, labelselector2)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				servicemonitor1: {},
				servicemonitor2: {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingServiceMonitorTestSuite) TestInvalidSelector() {
	s.AddServiceWithLabels(service1, &invalidlabelselector)
	s.AddServiceMonitorWithLabelSelector(servicemonitor1, invalidlabelselector)
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				servicemonitor1: {{Message: "service monitor has invalid label selector: key: Invalid value: \"-incorrect-labelselector-and-it-is-too-long-for-kubernetes-today-\": name part must be no more than 63 characters; name part must consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character (e.g. 'MyName',  or 'my.name',  or '123-abc', regex used for validation is '([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]')"}},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingServiceMonitorTestSuite) TestOneDanglingServiceMonitorIsDangling() {
	s.AddServiceWithLabels(service2, &labelselector2)
	s.AddServiceMonitorWithLabelSelector(servicemonitor1, labelselector1)
	s.AddServiceMonitorWithLabelSelector(servicemonitor2, labelselector2)
	label1, _ := metaV1.LabelSelectorAsSelector(&metaV1.LabelSelector{MatchLabels: labelselector1.MatchLabels})
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				servicemonitor1: {{Message: fmt.Sprintf("no services found matching the service monitor's label selector (%v) and namespace selector ([])", label1)}},
				servicemonitor2: {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingServiceMonitorTestSuite) TestNamespaceSelector() {
	s.AddServiceWithNamespace(service1, namespace1)
	s.AddServiceMonitorWithNamespaceSelector(servicemonitor1, namespaceselector1)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				servicemonitor1: {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingServiceMonitorTestSuite) TestDanglingNamespaceSelector() {
	s.AddServiceWithNamespace(service1, namespace2)
	s.AddServiceMonitorWithNamespaceSelector(servicemonitor1, namespaceselector1)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				servicemonitor1: {{Message: fmt.Sprintf("no services found matching the service monitor's label selector () and namespace selector (%v)", namespaceselector1.MatchNames)}},
			},
			ExpectInstantiationError: false,
		},
	})
}
