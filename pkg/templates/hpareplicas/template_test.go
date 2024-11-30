package hpareplicas

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/hpareplicas/internal/params"
	autoscalingV1 "k8s.io/api/autoscaling/v1"
	autoscalingV2 "k8s.io/api/autoscaling/v2"
	autoscalingV2Beta1 "k8s.io/api/autoscaling/v2beta1"
	autoscalingV2Beta2 "k8s.io/api/autoscaling/v2beta2"
)

var autoscalingVersions = [4]string{"v2beta1", "v2beta2", "v2", "v1"}

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

func (s *HPAReplicaTestSuite) addHPAWithReplicas(name string, replicas int32, version string) {
	s.ctx.AddMockHorizontalPodAutoscaler(s.T(), name, version)
	switch version {
	case "v2beta1":
		s.ctx.ModifyHorizontalPodAutoscalerV2Beta1(s.T(), name, func(hpa *autoscalingV2Beta1.HorizontalPodAutoscaler) {
			hpa.Spec.MinReplicas = &replicas
		})
	case "v2beta2":
		s.ctx.ModifyHorizontalPodAutoscalerV2Beta2(s.T(), name, func(hpa *autoscalingV2Beta2.HorizontalPodAutoscaler) {
			hpa.Spec.MinReplicas = &replicas
		})
	case "v2":
		s.ctx.ModifyHorizontalPodAutoscalerV2(s.T(), name, func(hpa *autoscalingV2.HorizontalPodAutoscaler) {
			hpa.Spec.MinReplicas = &replicas
		})
	case "v1":
		s.ctx.ModifyHorizontalPodAutoscalerV1(s.T(), name, func(hpa *autoscalingV1.HorizontalPodAutoscaler) {
			hpa.Spec.MinReplicas = &replicas
		})
	default:
		s.Require().FailNow(fmt.Sprintf("Unknown autoscaling version %s", version))
	}
}

func (s *HPAReplicaTestSuite) TestTooFewReplicas() {
	const (
		noExplicitReplicasHPAName = "hpa-no-explicit-replicas"
		twoReplicasHPAName        = "hpa-two-replicas"
	)

	for _, version := range autoscalingVersions {
		s.ctx.AddMockHorizontalPodAutoscaler(s.T(), noExplicitReplicasHPAName, version)
		s.addHPAWithReplicas(twoReplicasHPAName, 2, version)

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
}

func (s *HPAReplicaTestSuite) TestAcceptableReplicas() {
	const (
		acceptableReplicasHPAName = "hpa-acceptable-replicas"
	)

	for _, version := range autoscalingVersions {
		s.addHPAWithReplicas(acceptableReplicasHPAName, 3, version)

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
}
