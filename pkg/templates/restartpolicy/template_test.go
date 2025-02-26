package restartpolicy

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/restartpolicy/internal/params"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
)

func TestRestartPolicy(t *testing.T) {
	suite.Run(t, new(RestartPolicyTestSuite))
}

type RestartPolicyTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *RestartPolicyTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *RestartPolicyTestSuite) addDeploymentWithRestartPolicy(name string, policy coreV1.RestartPolicy) {
	s.ctx.AddMockDeployment(s.T(), name)
	s.ctx.ModifyDeployment(s.T(), name, func(deployment *appsV1.Deployment) {
		deployment.Spec.Template.Spec.RestartPolicy = policy
	})
}

func (s *RestartPolicyTestSuite) addDeploymentWithEmptyRestartPolicy(name string) {
	s.ctx.AddMockDeployment(s.T(), name)
	s.ctx.ModifyDeployment(s.T(), name, func(deployment *appsV1.Deployment) {
		deployment.Spec.Template.Spec.RestartPolicy = ""
	})
}

func (s *RestartPolicyTestSuite) addDeploymentWithoutRestartPolicy(name string) {
	s.ctx.AddMockDeployment(s.T(), name)
}

func (s *RestartPolicyTestSuite) addObjectWithoutPodSpec(name string) {
	s.ctx.AddMockService(s.T(), name)
}

func (s *RestartPolicyTestSuite) TestInvalidRestartPolicies() {
	const (
		withoutRestartPolicy = "without-restart-policy"
		emptyRestartPolicy   = "empty-restart-policy"
		restartPolicyNever   = "restart-policy-never"
	)

	s.addDeploymentWithoutRestartPolicy(withoutRestartPolicy)
	s.addDeploymentWithEmptyRestartPolicy(emptyRestartPolicy)
	s.addDeploymentWithRestartPolicy(restartPolicyNever, coreV1.RestartPolicyNever)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				withoutRestartPolicy: {
					{Message: "object has a restart policy defined with '' but the only accepted restart policies are '[Always OnFailure]'"},
				},
				emptyRestartPolicy: {
					{Message: "object has a restart policy defined with '' but the only accepted restart policies are '[Always OnFailure]'"},
				},
				restartPolicyNever: {
					{Message: "object has a restart policy defined with 'Never' but the only accepted restart policies are '[Always OnFailure]'"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *RestartPolicyTestSuite) TestAcceptableRestartPolicy() {
	const (
		alwaysRestartPolicy    = "restart-policy-always"
		onFailureRestartPolicy = "restart-policy-on-failure"
	)
	s.addDeploymentWithRestartPolicy(alwaysRestartPolicy, coreV1.RestartPolicyAlways)
	s.addDeploymentWithRestartPolicy(onFailureRestartPolicy, coreV1.RestartPolicyOnFailure)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				alwaysRestartPolicy:    nil,
				onFailureRestartPolicy: nil,
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *RestartPolicyTestSuite) TestObjectWithoutPodSpec() {
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
