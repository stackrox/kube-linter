package memoryrequirements

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/internal/pointers"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/memoryrequirements/internal/params"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestMemoryRequirements(t *testing.T) {
	suite.Run(t, new(MemoryRequirementsTestSuite))
}

type MemoryRequirementsTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *MemoryRequirementsTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func requirements(memoryMB int64) v1.ResourceRequirements {
	return v1.ResourceRequirements{
		Requests: v1.ResourceList{v1.ResourceMemory: *resource.NewQuantity(memoryMB*bytesInMB, resource.BinarySI)},
		Limits:   v1.ResourceList{v1.ResourceMemory: *resource.NewQuantity(memoryMB*bytesInMB, resource.BinarySI)},
	}
}

// TestInitAndRegularContainersWithLimitsSet is a regression test for
// https://github.com/stackrox/kube-linter/issues/428: a deployment where both
// init and regular containers have memory requirements set should not be flagged.
func (s *MemoryRequirementsTestSuite) TestInitAndRegularContainersWithLimitsSet() {
	const deploymentName = "dep-with-init-and-regular-containers"
	s.ctx.AddMockDeployment(s.T(), deploymentName)
	s.ctx.AddInitContainerToDeployment(s.T(), deploymentName, v1.Container{
		Name:      "init",
		Resources: requirements(128),
	})
	s.ctx.AddContainerToDeployment(s.T(), deploymentName, v1.Container{
		Name:      "app",
		Resources: requirements(128),
	})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				RequirementsType: "limit",
				UpperBoundMB:     pointers.Int(0),
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentName: nil,
			},
		},
	})
}

// TestInitContainerMissingLimitIsFlagged ensures that an init container without
// memory requirements is still flagged, even when a sibling regular container has
// its requirements set.
func (s *MemoryRequirementsTestSuite) TestInitContainerMissingLimitIsFlagged() {
	const deploymentName = "dep-with-init-container-missing-limit"
	s.ctx.AddMockDeployment(s.T(), deploymentName)
	s.ctx.AddInitContainerToDeployment(s.T(), deploymentName, v1.Container{Name: "init"})
	s.ctx.AddContainerToDeployment(s.T(), deploymentName, v1.Container{
		Name:      "app",
		Resources: requirements(128),
	})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				RequirementsType: "limit",
				UpperBoundMB:     pointers.Int(0),
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentName: {
					{Message: `container "init" has memory limit 0`},
				},
			},
		},
	})
}
