package priorityclassname

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/priorityclassname/internal/params"
	appsV1 "k8s.io/api/apps/v1"
)

func TestPriorityClassName(t *testing.T) {
	suite.Run(t, new(PriorityClassNameTestSuite))
}

type PriorityClassNameTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *PriorityClassNameTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *PriorityClassNameTestSuite) addDeploymentWithPriorityClassName(name, priorityClassName string) {
	s.ctx.AddMockDeployment(s.T(), name)
	s.ctx.ModifyDeployment(s.T(), name, func(deployment *appsV1.Deployment) {
		deployment.Spec.Template.Spec.PriorityClassName = priorityClassName
	})
}

func (s *PriorityClassNameTestSuite) addDeploymentWithEmptyPriorityClassName(name string) {
	s.ctx.AddMockDeployment(s.T(), name)
	s.ctx.ModifyDeployment(s.T(), name, func(deployment *appsV1.Deployment) {
		deployment.Spec.Template.Spec.PriorityClassName = ""
	})
}

func (s *PriorityClassNameTestSuite) addDeploymentWithoutPriorityClassName(name string) {
	s.ctx.AddMockDeployment(s.T(), name)
}

func (s *PriorityClassNameTestSuite) addObjectWithoutPodSpec(name string) {
	s.ctx.AddMockService(s.T(), name)
}

func (s *PriorityClassNameTestSuite) TestInvalidPriorityClassName() {
	const (
		invalidPriorityClassName = "invalid-priority-class-name"
	)

	s.addDeploymentWithPriorityClassName(invalidPriorityClassName, "system-node-critical")

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				AcceptedPriorityClassNames: []string{"system-cluster-critical", "custom-priority-class-name"},
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				invalidPriorityClassName: {
					{Message: "object has a priority class name defined with 'system-node-critical' but the only accepted priority class names are '[system-cluster-critical custom-priority-class-name]'"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *PriorityClassNameTestSuite) TestAcceptablePriorityClassName() {
	const (
		validPriorityClassName   = "valid-priority-class-name"
		emptyPriorityClassName   = "empty-priotity-class-name"
		withoutPriorityClassName = "without-piority-class-name"
	)

	s.addDeploymentWithPriorityClassName(validPriorityClassName, "system-cluster-critical")
	s.addDeploymentWithEmptyPriorityClassName(emptyPriorityClassName)
	s.addDeploymentWithoutPriorityClassName(withoutPriorityClassName)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				AcceptedPriorityClassNames: []string{"system-cluster-critical", "custom-priority-class-name"},
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				validPriorityClassName:   nil,
				emptyPriorityClassName:   nil,
				withoutPriorityClassName: nil,
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *PriorityClassNameTestSuite) TestObjectWithoutPodSpec() {
	const (
		objectWithoutPodSpec = "object-without-pod-spec"
	)

	s.addObjectWithoutPodSpec(objectWithoutPodSpec)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				objectWithoutPodSpec: nil,
			},
			ExpectInstantiationError: false,
		},
	})
}
