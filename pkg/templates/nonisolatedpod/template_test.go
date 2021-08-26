package nonisolatedpod

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/nonisolatedpod/internal/params"

	appsV1 "k8s.io/api/apps/v1"
	networkingV1 "k8s.io/api/networking/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	pod1             = "pod1"
	pod2             = "pod2"
	networkpolicyall = "networkpolicy-isolating-all"
	networkpolicy1   = "networkpolicy-isolating-pod1"
	networkpolicy2   = "networkpolicy-isolating-pod2"
	networkpolicy3   = "networkpolicy-isolating-none"
)

var emptyLabelSelector = metaV1.LabelSelector{} //empty selector

var labelselector1 = metaV1.LabelSelector{
	MatchLabels: map[string]string{"app": "pod1-test"},
}

var labelselector2 = metaV1.LabelSelector{
	MatchLabels: map[string]string{"app": "pod2-test"},
}

var labelselector3 = metaV1.LabelSelector{
	MatchLabels: map[string]string{"app": "no-pods-test"},
}

func TestNonIsolatedPod(t *testing.T) {
	suite.Run(t, new(NonIsolatedPodTestSuite))
}

type NonIsolatedPodTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *NonIsolatedPodTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *NonIsolatedPodTestSuite) AddNetworkPolicy(name string, podSelector metaV1.LabelSelector) {
	s.ctx.AddMockNetworkPolicy(s.T(), name)
	s.ctx.ModifyNetworkPolicy(s.T(), name, func(networkpolicy *networkingV1.NetworkPolicy) {
		networkpolicy.Spec.PodSelector = podSelector
	})
}

func (s *NonIsolatedPodTestSuite) AddDeploymentWithLabels(name string, labels *metaV1.LabelSelector) {
	s.ctx.AddMockDeployment(s.T(), name)
	s.ctx.ModifyDeployment(s.T(), name, func(deployment *appsV1.Deployment) {
		deployment.Spec.Template.Labels = labels.MatchLabels
	})
}

func (s *NonIsolatedPodTestSuite) TestAllPodsIsolatedWithNetworkPolicyEmptySelector() {
	s.AddDeploymentWithLabels(pod1, &labelselector1)
	s.AddDeploymentWithLabels(pod2, &labelselector2)
	s.AddNetworkPolicy(networkpolicyall, emptyLabelSelector)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				pod1: {},
				pod2: {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *NonIsolatedPodTestSuite) TestAllPodsIsolatedWithMultipleNetworkPolicies() {
	s.AddDeploymentWithLabels(pod1, &labelselector1)
	s.AddDeploymentWithLabels(pod2, &labelselector2)
	s.AddNetworkPolicy(networkpolicy1, labelselector1)
	s.AddNetworkPolicy(networkpolicy2, labelselector2)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				pod1: {},
				pod2: {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *NonIsolatedPodTestSuite) TestSomePodsIsolated() {
	s.AddDeploymentWithLabels(pod1, &labelselector1)
	s.AddDeploymentWithLabels(pod2, &labelselector2)
	s.AddNetworkPolicy(networkpolicy2, labelselector2)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				pod1: {{Message: "pods created by this object are non-isolated"}},
				pod2: {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *NonIsolatedPodTestSuite) TestAllPodsNonIsolated() {
	s.AddDeploymentWithLabels(pod1, &labelselector1)
	s.AddDeploymentWithLabels(pod2, &labelselector2)
	s.AddNetworkPolicy(networkpolicy3, labelselector3)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				pod1: {{Message: "pods created by this object are non-isolated"}},
				pod2: {{Message: "pods created by this object are non-isolated"}},
			},
			ExpectInstantiationError: false,
		},
	})
}
