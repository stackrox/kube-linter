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
	s.Init("verify-container-capabilities")
	s.ctx = mocks.NewMockContext()
}

func (s *ContainerCapabilitiesTestSuite) TestForbiddenCapabilities() {
	forbiddenCap := "NET_ADMIN"
	s.addPod()
	s.ctx.AddContainerToPod(podName, containerName, "", nil, nil, &v1.SecurityContext{
		Capabilities: &v1.Capabilities{
			Add:  []v1.Capability{v1.Capability(forbiddenCap)},
			Drop: nil,
		},
	})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				ForbiddenCapabilities: []string{forbiddenCap},
				Exceptions:            nil,
			},
			Diagnostics: []diagnostic.Diagnostic{
				{Message: fmt.Sprintf(addListDiagMsgFmt, containerName, forbiddenCap)},
				{Message: fmt.Sprintf(dropListDiagMsgFmt, containerName, []string{}, forbiddenCap)},
			},
			ExpectError: false,
		},
	})
}

func (s *ContainerCapabilitiesTestSuite) addPod() {
	s.ctx.AddMockPod(podName, "", "", nil, nil)
}