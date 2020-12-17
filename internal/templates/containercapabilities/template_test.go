package containercapabilities

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/internal/diagnostic"
	"golang.stackrox.io/kube-linter/internal/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/internal/templates"
	"golang.stackrox.io/kube-linter/internal/templates/containercapabilities/internal/params"
	v1 "k8s.io/api/core/v1"
)

var (
	podName       = "test-pod"
	containerName = "test-container"
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
	forbiddenCap := "SOME_CAP"
	addCaps := []v1.Capability{v1.Capability(forbiddenCap), "ALLOWED_CAP"}
	dropCaps := []v1.Capability{"DROPPED_CAP"}

	s.addPod()
	s.ctx.AddContainerToPod(podName, containerName, "", nil, nil, &v1.SecurityContext{
		Capabilities: &v1.Capabilities{
			Add:  addCaps,
			Drop: dropCaps,
		},
	})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				ForbiddenCapabilities: []string{forbiddenCap, "DROPPED_CAP"},
				Exceptions:            nil,
			},
			Diagnostics: []diagnostic.Diagnostic{
				{Message: fmt.Sprintf(addListDiagMsgFmt, containerName, forbiddenCap)},
				{Message: fmt.Sprintf(dropListDiagMsgFmt, containerName, dropCaps, forbiddenCap)},
			},
			ExpectError: false,
		},
	})
}

func (s *ContainerCapabilitiesTestSuite) TestForbiddenCapabilitiesWithAll() {

}

func (s *ContainerCapabilitiesTestSuite) addPod() {
	s.ctx.AddMockPod(podName, "", "", nil, nil)
}