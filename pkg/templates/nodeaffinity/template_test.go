package nodeaffinity

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/nodeaffinity/internal/params"
	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

func TestNodeAffinity(t *testing.T) {
	suite.Run(t, new(NodeAffinityTestSuite))
}

type NodeAffinityTestSuite struct {
	templates.TemplateTestSuite
	ctx *mocks.MockLintContext
}

const (
	templateKey              = "no-node-affinity"
	deploymentName           = "deployment"
	nodeAffinityErrorMessage = "object does not define any node affinity rules."
)

func (s *NodeAffinityTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *NodeAffinityTestSuite) TestIgnoreNodeAffinityCheckOnObjectWithoutAffinity() {
	s.ctx.AddMockClusterRole(s.T(), deploymentName)
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param:                    params.Params{},
			Diagnostics:              nil,
			ExpectInstantiationError: false,
		},
	})
}

func (s *NodeAffinityTestSuite) TestNoPodTemplateSpecAffinityDefined() {
	s.ctx.AddMockDeployment(s.T(), deploymentName)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentName: {
					{
						Message: nodeAffinityErrorMessage,
					},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *NodeAffinityTestSuite) TestNoNodeAffinityDefined() {
	s.ctx.AddMockDeployment(s.T(), deploymentName)

	s.ctx.ModifyDeployment(s.T(), deploymentName, func(deployment *appsV1.Deployment) {
		affinity := &v1.Affinity{
			NodeAffinity: nil,
		}
		deployment.Spec.Template.Spec.Affinity = affinity
	})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentName: {
					{
						Message: nodeAffinityErrorMessage,
					},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *NodeAffinityTestSuite) TestNodeAffinityDefined() {
	s.ctx.AddMockDeployment(s.T(), deploymentName)

	s.ctx.ModifyDeployment(s.T(), deploymentName, func(deployment *appsV1.Deployment) {
		affinity := &v1.Affinity{
			NodeAffinity: &v1.NodeAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: &v1.NodeSelector{
					NodeSelectorTerms: []v1.NodeSelectorTerm{
						{
							MatchExpressions: []v1.NodeSelectorRequirement{
								{
									Key:      "nodeKey",
									Operator: "In",
									Values: []string{
										"NodeA",
										"NodeB",
										"NodeC",
									},
								},
							},
						},
					},
				},
			},
		}
		deployment.Spec.Template.Spec.Affinity = affinity
	})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param:                    params.Params{},
			Diagnostics:              nil,
			ExpectInstantiationError: false,
		},
	})
}
