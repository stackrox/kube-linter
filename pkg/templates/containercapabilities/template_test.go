package containercapabilities

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/containercapabilities/internal/params"
	v1 "k8s.io/api/core/v1"
)

var (
	deploymentName = "test-deployment"
	containerName  = "test-container"
)

func TestContainerCapabilities(t *testing.T) {
	suite.Run(t, new(ContainerCapabilitiesTestSuite))
}

type ContainerCapabilitiesTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *ContainerCapabilitiesTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *ContainerCapabilitiesTestSuite) TestForbiddenCapabilities() {
	addCaps := []v1.Capability{"FORBIDDEN_CAP", "ALLOWED_CAP"}
	dropCaps := []v1.Capability{"DROPPED_CAP"}

	s.addPodAndAddContainerWithCaps(addCaps, dropCaps)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				ForbiddenCapabilities: []string{"FORBIDDEN_CAP", "DROPPED_CAP"},
				Exceptions:            nil,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentName: {
					{Message: fmt.Sprintf(addListDiagMsgFmt, containerName, "FORBIDDEN_CAP")},
					{Message: fmt.Sprintf(dropListDiagMsgFmt, containerName, dropCaps, "FORBIDDEN_CAP")},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *ContainerCapabilitiesTestSuite) TestForbiddenCapabilitiesWithAll() {
	addCaps := []v1.Capability{"CAP_1", "CAP_2", "CAP_3"}
	dropCaps := []v1.Capability{"DROPPED_CAP"}

	s.addPodAndAddContainerWithCaps(addCaps, dropCaps)

	s.Validate(s.ctx, []templates.TestCase{
		// Case 1: all are prohibited
		{
			Param: params.Params{
				ForbiddenCapabilities: []string{"all"},
				Exceptions:            nil,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentName: {
					{Message: fmt.Sprintf(addListWithAllDiagMsgFmt, containerName, "CAP_1")},
					{Message: fmt.Sprintf(addListWithAllDiagMsgFmt, containerName, "CAP_2")},
					{Message: fmt.Sprintf(addListWithAllDiagMsgFmt, containerName, "CAP_3")},
					{Message: fmt.Sprintf(dropListWithAllDiagMsgFmt, containerName, dropCaps)},
				},
			},
			ExpectInstantiationError: false,
		},
		// Case 2: with some forgiven capabilities
		{
			Param: params.Params{
				// Also tests reserved word "all" should match irrespective of case
				ForbiddenCapabilities: []string{"AlL"},
				Exceptions:            []string{"CAP_1", "CAP_2"},
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentName: {
					{Message: fmt.Sprintf(addListWithAllDiagMsgFmt, containerName, "CAP_3")},
					{Message: fmt.Sprintf(dropListWithAllDiagMsgFmt, containerName, dropCaps)},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *ContainerCapabilitiesTestSuite) TestAddListHasAll() {
	addCaps := []v1.Capability{"all", "REDUNDANT_CAP"}
	dropCaps := make([]v1.Capability, 0)

	s.addPodAndAddContainerWithCaps(addCaps, dropCaps)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				ForbiddenCapabilities: []string{"CAP_1"},
				Exceptions:            nil,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentName: {
					{Message: fmt.Sprintf(addListDiagMsgFmt, containerName, "all")},
					{Message: fmt.Sprintf(dropListDiagMsgFmt, containerName, dropCaps, "CAP_1")},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *ContainerCapabilitiesTestSuite) TestDropListHasAll() {
	addCaps := []v1.Capability{"FORGIVEN_CAP"}
	dropCaps := []v1.Capability{"all"}

	s.addPodAndAddContainerWithCaps(addCaps, dropCaps)

	s.Validate(s.ctx, []templates.TestCase{
		// Case 1: caps are all dropped by "all" in drop list
		{
			Param: params.Params{
				ForbiddenCapabilities: []string{"CAP_1", "CAP_2"},
				Exceptions:            nil,
			},
			Diagnostics:              nil,
			ExpectInstantiationError: false,
		},
		// Case 2: forbidden caps include "all"
		{
			Param: params.Params{
				ForbiddenCapabilities: []string{"all"},
				Exceptions:            nil,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentName: {
					{Message: fmt.Sprintf(addListWithAllDiagMsgFmt, containerName, "FORGIVEN_CAP")},
				},
			},
			ExpectInstantiationError: false,
		},
		// Case 3: now we forgive the FORGIVEN_CAP. Should see no error
		{
			Param: params.Params{
				ForbiddenCapabilities: []string{"all"},
				Exceptions:            []string{"FORGIVEN_CAP"},
			},
			Diagnostics:              nil,
			ExpectInstantiationError: false,
		},
	})
}

func (s *ContainerCapabilitiesTestSuite) TestInvalidParams() {
	addCaps := []v1.Capability{"CAP_1"}
	dropCaps := []v1.Capability{"CAP_2"}

	s.addPodAndAddContainerWithCaps(addCaps, dropCaps)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				ForbiddenCapabilities: []string{"THIS_IS_NOT_All_CAP"},
				Exceptions:            []string{"BUT_WE_SPECIFY_EXCEPTIONS"},
			},
			Diagnostics:              nil,
			ExpectInstantiationError: true,
		},
	})
}

func (s *ContainerCapabilitiesTestSuite) addPodAndAddContainerWithCaps(addCaps, dropCaps []v1.Capability) {
	s.ctx.AddMockDeployment(s.T(), deploymentName)
	s.ctx.AddContainerToDeployment(s.T(), deploymentName, v1.Container{Name: containerName, SecurityContext: &v1.SecurityContext{
		Capabilities: &v1.Capabilities{
			Add:  addCaps,
			Drop: dropCaps,
		},
	}})
}
