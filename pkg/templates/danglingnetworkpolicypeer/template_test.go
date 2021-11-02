package danglingnetworkpolicypeer

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/danglingnetworkpolicypeer/internal/params"

	appsV1 "k8s.io/api/apps/v1"
	networkingV1 "k8s.io/api/networking/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	pod1                    = "pod1"
	pod2                    = "pod2"
	networkpolicyIngressAll = "networkpolicypeer-ingress-matches-all" // empty selector
	networkpolicyEgressAll  = "networkpolicypeer-egress-matches-all"
	networkpolicy1          = "networkpolicypeer-matches-pod1"
	networkpolicy2          = "networkpolicyeer-matches-pod2"
	networkpolicy3          = "networkpolicy-egress-multiple-peers"
	networkpolicy4          = "networkpolicy-ingress-egress-multiple-peers"
)

var emptyLabelSelector = metaV1.LabelSelector{} //empty selector

var labelselector1 = metaV1.LabelSelector{
	MatchLabels: map[string]string{"app": "pod1-test"},
}

var labelselector2 = metaV1.LabelSelector{
	MatchLabels: map[string]string{"app": "pod2-test"},
}

var labelselector3 = metaV1.LabelSelector{
	MatchLabels: map[string]string{"app": "wrong1-test"},
}

var labelselector4 = metaV1.LabelSelector{
	MatchLabels: map[string]string{"app": "wrong2-test"},
}

func TestDanglingNetworkPolicyPeer(t *testing.T) {
	suite.Run(t, new(DanglingNetworkPolicyPeerTestSuite))
}

type DanglingNetworkPolicyPeerTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *DanglingNetworkPolicyPeerTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *DanglingNetworkPolicyPeerTestSuite) AddNetworkPolicyWithEgress(name string, podSelector metaV1.LabelSelector) {
	egressPeer := networkingV1.NetworkPolicyPeer{
		PodSelector: &podSelector,
	}

	s.ctx.AddMockNetworkPolicy(s.T(), name)
	s.ctx.ModifyNetworkPolicy(s.T(), name, func(networkpolicy *networkingV1.NetworkPolicy) {
		networkpolicy.Spec.PodSelector = emptyLabelSelector
		networkpolicy.Spec.Egress = append(networkpolicy.Spec.Egress, networkingV1.NetworkPolicyEgressRule{
			To: []networkingV1.NetworkPolicyPeer{egressPeer},
		})
	})
}

func (s *DanglingNetworkPolicyPeerTestSuite) AddNetworkPolicyWithEgressMultiplePeers(name string, podSelector1, podSelector2, podSelector3 metaV1.LabelSelector) {
	egressPeer1 := networkingV1.NetworkPolicyPeer{
		PodSelector: &podSelector1,
	}
	egressPeer2 := networkingV1.NetworkPolicyPeer{
		PodSelector: &podSelector2,
	}

	egressPeer3 := networkingV1.NetworkPolicyPeer{
		PodSelector: &podSelector3,
	}

	s.ctx.AddMockNetworkPolicy(s.T(), name)
	s.ctx.ModifyNetworkPolicy(s.T(), name, func(networkpolicy *networkingV1.NetworkPolicy) {
		networkpolicy.Spec.PodSelector = emptyLabelSelector
		networkpolicy.Spec.Egress = append(networkpolicy.Spec.Egress, networkingV1.NetworkPolicyEgressRule{
			To: []networkingV1.NetworkPolicyPeer{egressPeer1, egressPeer2, egressPeer3},
		})
	})
}

func (s *DanglingNetworkPolicyPeerTestSuite) AddNetworkPolicyWithIngress(name string, podSelector metaV1.LabelSelector) {
	ingressPeer := networkingV1.NetworkPolicyPeer{
		PodSelector: &podSelector,
	}

	s.ctx.AddMockNetworkPolicy(s.T(), name)
	s.ctx.ModifyNetworkPolicy(s.T(), name, func(networkpolicy *networkingV1.NetworkPolicy) {
		networkpolicy.Spec.PodSelector = emptyLabelSelector
		networkpolicy.Spec.Ingress = append(networkpolicy.Spec.Ingress, networkingV1.NetworkPolicyIngressRule{
			From: []networkingV1.NetworkPolicyPeer{ingressPeer},
		})
	})
}

func (s *DanglingNetworkPolicyPeerTestSuite) AddNetworkPolicyWithIngressEgressMultiplePeers(name string, podSelector1, podSelector2, podSelector3, podSelector4 metaV1.LabelSelector) {
	ingressPeer1 := networkingV1.NetworkPolicyPeer{
		PodSelector: &podSelector1,
	}
	ingressPeer2 := networkingV1.NetworkPolicyPeer{
		PodSelector: &podSelector3,
	}
	egressPeer1 := networkingV1.NetworkPolicyPeer{
		PodSelector: &podSelector2,
	}
	egressPeer2 := networkingV1.NetworkPolicyPeer{
		PodSelector: &podSelector4,
	}

	s.ctx.AddMockNetworkPolicy(s.T(), name)
	s.ctx.ModifyNetworkPolicy(s.T(), name, func(networkpolicy *networkingV1.NetworkPolicy) {
		networkpolicy.Spec.PodSelector = emptyLabelSelector
		networkpolicy.Spec.Ingress = append(networkpolicy.Spec.Ingress, networkingV1.NetworkPolicyIngressRule{
			From: []networkingV1.NetworkPolicyPeer{ingressPeer1, ingressPeer2},
		})
		networkpolicy.Spec.Egress = append(networkpolicy.Spec.Egress, networkingV1.NetworkPolicyEgressRule{
			To: []networkingV1.NetworkPolicyPeer{egressPeer1, egressPeer2},
		})
	})
}

func (s *DanglingNetworkPolicyPeerTestSuite) AddDeploymentWithLabels(name string, labels *metaV1.LabelSelector) {
	s.ctx.AddMockDeployment(s.T(), name)
	s.ctx.ModifyDeployment(s.T(), name, func(deployment *appsV1.Deployment) {
		deployment.Spec.Template.Labels = labels.MatchLabels
	})
}

func (s *DanglingNetworkPolicyPeerTestSuite) TestNetworkPolicyPeerEmptySelectorMatchesAllPods() {
	s.AddDeploymentWithLabels(pod1, &labelselector1)
	s.AddDeploymentWithLabels(pod2, &labelselector2)
	s.AddNetworkPolicyWithEgress(networkpolicyEgressAll, emptyLabelSelector)
	s.AddNetworkPolicyWithIngress(networkpolicyIngressAll, emptyLabelSelector)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				networkpolicyEgressAll:  {},
				networkpolicyIngressAll: {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingNetworkPolicyPeerTestSuite) TestNoDanglingNetworkPolicyPeers() {
	s.AddDeploymentWithLabels(pod1, &labelselector1)
	s.AddDeploymentWithLabels(pod2, &labelselector2)
	s.AddNetworkPolicyWithEgress(networkpolicy1, labelselector1)
	s.AddNetworkPolicyWithIngress(networkpolicy2, labelselector2)

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

func (s *DanglingNetworkPolicyPeerTestSuite) TestOneNetworkPolicyPeerIsDangling() {
	s.AddDeploymentWithLabels(pod2, &labelselector2)
	s.AddNetworkPolicyWithEgress(networkpolicy1, labelselector1)
	s.AddNetworkPolicyWithEgress(networkpolicy2, labelselector2)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				networkpolicy1: {{Message: "no pods found matching networkpolicy rule's podSelector labels (app=pod1-test)"}},
				networkpolicy2: {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DanglingNetworkPolicyPeerTestSuite) TestEgressRuleWithMultiplePeerSomeDangling() {
	s.AddDeploymentWithLabels(pod1, &labelselector1)
	s.AddDeploymentWithLabels(pod2, &labelselector2)
	s.AddNetworkPolicyWithEgressMultiplePeers(networkpolicy3, labelselector1, labelselector3, labelselector4)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				networkpolicy3: {{Message: "no pods found matching networkpolicy rule's podSelector labels (app=wrong1-test)"},
					{Message: "no pods found matching networkpolicy rule's podSelector labels (app=wrong2-test)"}},
			},
			ExpectInstantiationError: false,
		},
	})

}

func (s *DanglingNetworkPolicyPeerTestSuite) TestNPwithIngressEgressRulesWithMultiplePeerSomeDangling() {
	s.AddDeploymentWithLabels(pod1, &labelselector1)
	s.AddDeploymentWithLabels(pod2, &labelselector2)
	s.AddNetworkPolicyWithIngressEgressMultiplePeers(networkpolicy4, labelselector1, labelselector2, labelselector3, labelselector4)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				networkpolicy4: {{Message: "no pods found matching networkpolicy rule's podSelector labels (app=wrong1-test)"},
					{Message: "no pods found matching networkpolicy rule's podSelector labels (app=wrong2-test)"}},
			},
			ExpectInstantiationError: false,
		},
	})

}
