package serviceaccount

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/internal/pointers"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/serviceaccount/internal/params"
	appsV1 "k8s.io/api/apps/v1"
)

func TestServiceAccount(t *testing.T) {
	suite.Run(t, new(ServiceAccountTestSuite))
}

type ServiceAccountTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *ServiceAccountTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *ServiceAccountTestSuite) addDeploymentWithServiceAccount(name, serviceAccountName string, automountServiceAccountToken *bool) {
	s.ctx.AddMockDeployment(s.T(), name)
	s.ctx.ModifyDeployment(s.T(), name, func(deployment *appsV1.Deployment) {
		deployment.Spec.Template.Spec.ServiceAccountName = serviceAccountName
		deployment.Spec.Template.Spec.AutomountServiceAccountToken = automountServiceAccountToken
	})
}

func (s *ServiceAccountTestSuite) TestServiceAccountName() {
	const (
		nonDefaultServiceAccount       = "non-default"
		noAutoMountServiceAccountToken = "no-auto-mount-sa-token"
	)
	s.addDeploymentWithServiceAccount(nonDefaultServiceAccount, nonDefaultServiceAccount, pointers.Bool(true))
	s.addDeploymentWithServiceAccount(noAutoMountServiceAccountToken, "", pointers.Bool(false))

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				ServiceAccount: "non-default",
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				nonDefaultServiceAccount: {
					{Message: "found matching serviceAccount (\"non-default\")"},
				},
				noAutoMountServiceAccountToken: {},
			},
			ExpectInstantiationError: false,
		},
		{
			Param: params.Params{
				ServiceAccount: "[a)", // Wrong Regex which should raise an error
			},
			ExpectInstantiationError: true,
		},
	})
}
