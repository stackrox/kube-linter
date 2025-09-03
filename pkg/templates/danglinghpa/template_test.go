package danglinghpa

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/danglinghpa/internal/params"
	appsV1 "k8s.io/api/apps/v1"
	autoscalingV1 "k8s.io/api/autoscaling/v1"
	autoscalingV2 "k8s.io/api/autoscaling/v2"
	autoscalingV2Beta1 "k8s.io/api/autoscaling/v2beta1"
	autoscalingV2Beta2 "k8s.io/api/autoscaling/v2beta2"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	targetKind       = "Deployment"
	targetName       = "deployment-1"
	targetAPIVersion = "apps/v1"
)

var (
	autoscalingVersions = [4]string{"v2beta1", "v2beta2", "v2", "v1"}
	target              = objectReference{
		Kind:       targetKind,
		Name:       targetName,
		APIVersion: targetAPIVersion,
	}
)

func TestDanglingHpa(t *testing.T) {
	suite.Run(t, new(DanglingHpaSuite))
}

type DanglingHpaSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *DanglingHpaSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *DanglingHpaSuite) addHpaWithTarget(name string, target objectReference, version string) {
	s.ctx.AddMockHorizontalPodAutoscaler(s.T(), name, version)
	switch version {
	case "v2beta1":
		s.ctx.ModifyHorizontalPodAutoscalerV2Beta1(s.T(), name, func(hpa *autoscalingV2Beta1.HorizontalPodAutoscaler) {
			hpa.Spec.ScaleTargetRef = autoscalingV2Beta1.CrossVersionObjectReference{
				Kind:       target.Kind,
				Name:       target.Name,
				APIVersion: target.APIVersion,
			}
		})
	case "v2beta2":
		s.ctx.ModifyHorizontalPodAutoscalerV2Beta2(s.T(), name, func(hpa *autoscalingV2Beta2.HorizontalPodAutoscaler) {
			hpa.Spec.ScaleTargetRef = autoscalingV2Beta2.CrossVersionObjectReference{
				Kind:       target.Kind,
				Name:       target.Name,
				APIVersion: target.APIVersion,
			}
		})
	case "v2":
		s.ctx.ModifyHorizontalPodAutoscalerV2(s.T(), name, func(hpa *autoscalingV2.HorizontalPodAutoscaler) {
			hpa.Spec.ScaleTargetRef = autoscalingV2.CrossVersionObjectReference{
				Kind:       target.Kind,
				Name:       target.Name,
				APIVersion: target.APIVersion,
			}
		})
	case "v1":
		s.ctx.ModifyHorizontalPodAutoscalerV1(s.T(), name, func(hpa *autoscalingV1.HorizontalPodAutoscaler) {
			hpa.Spec.ScaleTargetRef = autoscalingV1.CrossVersionObjectReference{
				Kind:       target.Kind,
				Name:       target.Name,
				APIVersion: target.APIVersion,
			}
		})
	default:
		s.Require().FailNow(fmt.Sprintf("Unknown autoscaling version %s", version))
	}
}

func (s *DanglingHpaSuite) TestDanglingHpa() {
	const (
		noExplicitDeployment = "hpa-no-explicit-deployment"
		missingDeployment    = "hpa-missing-deployment"
	)

	for _, version := range autoscalingVersions {
		s.ctx.AddMockHorizontalPodAutoscaler(s.T(), noExplicitDeployment, version)
		s.addHpaWithTarget(missingDeployment, target, version)

		s.Validate(s.ctx, []templates.TestCase{
			{
				Param: params.Params{},
				Diagnostics: map[string][]diagnostic.Diagnostic{
					noExplicitDeployment: {
						{Message: "no resources found matching HorizontalPodAutoscaler scaleTargetRef ({  })"},
					},
					missingDeployment: {
						{Message: "no resources found matching HorizontalPodAutoscaler scaleTargetRef ({Deployment deployment-1 apps/v1})"},
					},
				},
				ExpectInstantiationError: false,
			},
		})
	}
}

func (s *DanglingHpaSuite) TestHpaPassesWhenDeploymentIsPresent() {
	const (
		nonDanglingHpa = "hpa-with-deployment"
	)
	s.ctx.AddMockDeployment(s.T(), targetName)
	s.ctx.ModifyDeployment(s.T(), targetName, func(deployment *appsV1.Deployment) {
		deployment.TypeMeta = metaV1.TypeMeta{
			Kind:       targetKind,
			APIVersion: targetAPIVersion,
		}
	})

	for _, version := range autoscalingVersions {
		s.addHpaWithTarget(nonDanglingHpa, target, version)

		s.Validate(s.ctx, []templates.TestCase{
			{
				Param: params.Params{},
				Diagnostics: map[string][]diagnostic.Diagnostic{
					nonDanglingHpa: nil,
				},
				ExpectInstantiationError: false,
			},
		})
	}
}
