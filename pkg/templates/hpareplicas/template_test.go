package hpareplicas

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/hpareplicas/internal/params"
	autoscalingV2Beta1 "k8s.io/api/autoscaling/v2beta1"
)

func TestHPAReplicas(t *testing.T) {
	suite.Run(t, new(HPAReplicaTestSuite))
}

type HPAReplicaTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *HPAReplicaTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *HPAReplicaTestSuite) addHPAWithReplicas(name string, replicas int32) {
	s.ctx.AddMockHorizontalPodAutoscaler(s.T(), name)
	s.ctx.ModifyHorizontalPodAutoscaler(s.T(), name, func(hpa *autoscalingV2Beta1.HorizontalPodAutoscaler) {
		hpa.Spec.MinReplicas = &replicas
	})
}

func (s *HPAReplicaTestSuite) TestTooFewReplicas() {
	const (
		noExplicitReplicasHPAName = "hpa-no-explicit-replicas"
		twoReplicasHPAName        = "hpa-two-replicas"
	)

	s.ctx.AddMockHorizontalPodAutoscaler(s.T(), noExplicitReplicasHPAName)
	s.addHPAWithReplicas(twoReplicasHPAName, 2)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				MinReplicas: 3,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				noExplicitReplicasHPAName: {
					{Message: "object has 1 replica but minimum required replicas is 3"},
				},
				twoReplicasHPAName: {
					{Message: "object has 2 replicas but minimum required replicas is 3"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *HPAReplicaTestSuite) TestAcceptableReplicas() {
	const (
		acceptableReplicasHPAName = "hpa-acceptable-replicas"
	)
	s.addHPAWithReplicas(acceptableReplicasHPAName, 3)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				MinReplicas: 3,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				acceptableReplicasHPAName: nil,
			},
			ExpectInstantiationError: false,
		},
	})
}
