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
		matchingSAWithAutoMountTokenNil      = "match-sa-token-mount-nil"
		matchingSAWithAutoMountTokenTrue     = "match-sa-token-mount-true"
		matchingSAWithAutoMountTokenFalse    = "match-sa-token-mount-false"
		nonMatchingSAWithAutoMountTokenNil   = "non-match-sa-token-mount-nil"
		nonMatchingSAWithAutoMountTokenTrue  = "non-match-sa-token-mount-true"
		nonMatchingSAWithAutoMountTokenFalse = "non-match-sa-token-mount-false"
	)

	s.addDeploymentWithServiceAccount(matchingSAWithAutoMountTokenNil, "non-default", nil)
	s.addDeploymentWithServiceAccount(matchingSAWithAutoMountTokenTrue, "non-default", pointers.Bool(true))
	s.addDeploymentWithServiceAccount(matchingSAWithAutoMountTokenFalse, "non-default", pointers.Bool(false))
	s.addDeploymentWithServiceAccount(nonMatchingSAWithAutoMountTokenNil, "", nil)
	s.addDeploymentWithServiceAccount(nonMatchingSAWithAutoMountTokenTrue, "", pointers.Bool(true))
	s.addDeploymentWithServiceAccount(nonMatchingSAWithAutoMountTokenFalse, "", pointers.Bool(false))

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				ServiceAccount: "non-default",
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				matchingSAWithAutoMountTokenNil: {
					{Message: "found matching serviceAccount (\"non-default\")"},
				},
				matchingSAWithAutoMountTokenTrue: {
					{Message: "found matching serviceAccount (\"non-default\")"},
				},
				matchingSAWithAutoMountTokenFalse:    {},
				nonMatchingSAWithAutoMountTokenNil:   {},
				nonMatchingSAWithAutoMountTokenTrue:  {},
				nonMatchingSAWithAutoMountTokenFalse: {},
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
