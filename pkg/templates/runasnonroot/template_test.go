package runasnonroot

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/runasnonroot/internal/params"
	v1 "k8s.io/api/core/v1"
)

func TestRunAsNonRoot(t *testing.T) {
	suite.Run(t, new(RunAsNonRootTestSuite))
}

type RunAsNonRootTestSuite struct {
	templates.TemplateTestSuite
	ctx *mocks.MockLintContext
}

func (s *RunAsNonRootTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *RunAsNonRootTestSuite) TestContainerRunAsGroupZero() {
	const deploymentName = "container-group-zero"

	s.ctx.AddMockDeployment(s.T(), deploymentName)
	s.ctx.AddContainerToDeployment(s.T(), deploymentName, v1.Container{
		Name: "app",
		SecurityContext: &v1.SecurityContext{
			RunAsUser:  int64Ptr(1000),
			RunAsGroup: int64Ptr(0),
		},
	})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentName: {
					{Message: `container "app" is set to runAsGroup 0`},
				},
			},
		},
	})
}

func (s *RunAsNonRootTestSuite) TestPodRunAsGroupZero() {
	const deploymentName = "pod-group-zero"

	s.ctx.AddMockDeployment(s.T(), deploymentName)
	s.ctx.AddSecurityContextToDeployment(s.T(), deploymentName, &v1.PodSecurityContext{
		RunAsGroup: int64Ptr(0),
	})
	s.ctx.AddContainerToDeployment(s.T(), deploymentName, v1.Container{
		Name: "app",
		SecurityContext: &v1.SecurityContext{
			RunAsUser: int64Ptr(1000),
		},
	})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentName: {
					{Message: `container "app" is set to runAsGroup 0`},
				},
			},
		},
	})
}

func (s *RunAsNonRootTestSuite) TestContainerRunAsGroupOverridesPodRunAsGroup() {
	const deploymentName = "container-group-overrides-pod"

	s.ctx.AddMockDeployment(s.T(), deploymentName)
	s.ctx.AddSecurityContextToDeployment(s.T(), deploymentName, &v1.PodSecurityContext{
		RunAsGroup: int64Ptr(0),
	})
	s.ctx.AddContainerToDeployment(s.T(), deploymentName, v1.Container{
		Name: "app",
		SecurityContext: &v1.SecurityContext{
			RunAsUser:  int64Ptr(1000),
			RunAsGroup: int64Ptr(1000),
		},
	})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param:       params.Params{},
			Diagnostics: nil,
		},
	})
}

func (s *RunAsNonRootTestSuite) TestMissingRunAsGroupAllowedWhenRunAsUserNonRoot() {
	const deploymentName = "non-root-user-no-group"

	s.ctx.AddMockDeployment(s.T(), deploymentName)
	s.ctx.AddContainerToDeployment(s.T(), deploymentName, v1.Container{
		Name: "app",
		SecurityContext: &v1.SecurityContext{
			RunAsUser: int64Ptr(1000),
		},
	})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param:       params.Params{},
			Diagnostics: nil,
		},
	})
}

func int64Ptr(v int64) *int64 {
	return &v
}
