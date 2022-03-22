package danglingnetworkpolicy

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/danglingnetworkpolicy/internal/params"

	appsV1 "k8s.io/api/apps/v1"
	networkingV1 "k8s.io/api/networking/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	pod1             = "pod1"
	pod2             = "pod2"
	networkpolicyall = "networkpolicy-matches-all" // empty selector
	networkpolicy1   = "networkpolicy-matches-pod1"
	networkpolicy2   = "networkpolicy-matches-pod2"
)

var emptyLabelSelector = metaV1.LabelSelector{} //empty selector

var labelselector1 = metaV1.LabelSelector{
	MatchLabels: map[string]string{"app": "pod1-test"},
}

var labelselector2 = metaV1.LabelSelector{
	MatchLabels: map[string]string{"app": "pod2-test"},
}

func TestDanglingNetworkPolicy(t *testing.T) {
	suite.Run(t, new(DanglingNetworkPolicyTestSuite))
}

type DanglingNetworkPolicyTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *DanglingNetworkPolicyTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *DanglingNetworkPolicyTestSuite) AddNetworkPolicy(name string, podSelector metaV1.LabelSelector) {
	s.ctx.AddMockNetworkPolicy(s.T(), name)
	s.ctx.ModifyNetworkPolicy(s.T(), name, func(networkpolicy *networkingV1.NetworkPolicy) {
		networkpolicy.Spec.PodSelector = podSelector
	})
}

func (s *DanglingNetworkPolicyTestSuite) AddDeploymentWithLabels(name string, labels *metaV1.LabelSelector) {
	s.ctx.AddMockDeployment(s.T(), name)
	s.ctx.ModifyDeployment(s.T(), name, func(deployment *appsV1.Deployment) {
		deployment.Spec.Template.Labels = labels.MatchLabels
	})
}

func (s *DanglingNetworkPolicyTestSuite) TestNetworkPolicyEmptySelectorMatchesAllPods() {
	s.AddDeploymentWithLabels(pod1, &labelselector1)
	s.AddDeploymentWithLabels(pod2, &labelselector2)
	s.AddNetworkPolicy(networkpolicyall, emptyLabelSelector)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				networkpolicyall: {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingNetworkPolicyTestSuite) TestNoDanglingNetworkPolicies() {
	s.AddDeploymentWithLabels(pod1, &labelselector1)
	s.AddDeploymentWithLabels(pod2, &labelselector2)
	s.AddNetworkPolicy(networkpolicy1, labelselector1)
	s.AddNetworkPolicy(networkpolicy2, labelselector2)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				networkpolicy1: {},
				networkpolicy2: {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingNetworkPolicyTestSuite) TestOneNetworkPolicyIsDangling() {
	s.AddDeploymentWithLabels(pod2, &labelselector2)
	s.AddNetworkPolicy(networkpolicy1, labelselector1)
	s.AddNetworkPolicy(networkpolicy2, labelselector2)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				networkpolicy1: {{Message: fmt.Sprintf("no pods found matching networkpolicy's podSelector labels (%v) ", labelselector1)}},
				networkpolicy2: {},
			},
			ExpectInstantiationError: false,
		},
	})
}
