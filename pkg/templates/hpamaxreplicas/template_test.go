package hpamaxreplicas

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/hpamaxreplicas/internal/params"
	autoscalingV1 "k8s.io/api/autoscaling/v1"
	autoscalingV2 "k8s.io/api/autoscaling/v2"
	autoscalingV2Beta1 "k8s.io/api/autoscaling/v2beta1"
	autoscalingV2Beta2 "k8s.io/api/autoscaling/v2beta2"
)

var autoscalingVersions = [4]string{"v2beta1", "v2beta2", "v2", "v1"}

func TestHPAMaxReplicas(t *testing.T) {
	suite.Run(t, new(HPAMaxReplicaTestSuite))
}

type HPAMaxReplicaTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *HPAMaxReplicaTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *HPAMaxReplicaTestSuite) addHPAWithMaxReplicas(name string, replicas int32, version string) {
	s.ctx.AddMockHorizontalPodAutoscaler(s.T(), name, version)
	switch version {
	case "v2beta1":
		s.ctx.ModifyHorizontalPodAutoscalerV2Beta1(s.T(), name, func(hpa *autoscalingV2Beta1.HorizontalPodAutoscaler) {
			hpa.Spec.MaxReplicas = replicas
		})
	case "v2beta2":
		s.ctx.ModifyHorizontalPodAutoscalerV2Beta2(s.T(), name, func(hpa *autoscalingV2Beta2.HorizontalPodAutoscaler) {
			hpa.Spec.MaxReplicas = replicas
		})
	case "v2":
		s.ctx.ModifyHorizontalPodAutoscalerV2(s.T(), name, func(hpa *autoscalingV2.HorizontalPodAutoscaler) {
			hpa.Spec.MaxReplicas = replicas
		})
	case "v1":
		s.ctx.ModifyHorizontalPodAutoscalerV1(s.T(), name, func(hpa *autoscalingV1.HorizontalPodAutoscaler) {
			hpa.Spec.MaxReplicas = replicas
		})
	default:
		s.Require().FailNow(fmt.Sprintf("Unknown autoscaling version %s", version))
	}
}

func (s *HPAMaxReplicaTestSuite) TestTooManyReplicas() {
	const (
		tenReplicasHPAName = "hpa-ten-replicas"
	)

	for _, version := range autoscalingVersions {
		s.addHPAWithMaxReplicas(tenReplicasHPAName, 10, version)

		s.Validate(s.ctx, []templates.TestCase{
			{
				Param: params.Params{
					MaxReplicas: 5,
				},
				Diagnostics: map[string][]diagnostic.Diagnostic{
					tenReplicasHPAName: {
						{Message: "object has 10 replicas but maximum allowed replicas is 5"},
					},
				},
				ExpectInstantiationError: false,
			},
		})
	}
}

func (s *HPAMaxReplicaTestSuite) TestAcceptableReplicas() {
	const (
		acceptableReplicasHPAName = "hpa-acceptable-replicas"
	)

	for _, version := range autoscalingVersions {
		s.addHPAWithMaxReplicas(acceptableReplicasHPAName, 5, version)

		s.Validate(s.ctx, []templates.TestCase{
			{
				Param: params.Params{
					MaxReplicas: 5,
				},
				Diagnostics: map[string][]diagnostic.Diagnostic{
					acceptableReplicasHPAName: nil,
				},
				ExpectInstantiationError: false,
			},
		})
	}
}
