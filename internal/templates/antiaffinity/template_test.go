package antiaffinity

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

func TestAntiAffinity(t *testing.T) {
	suite.Run(t, new(AntiAffinityTestSuite))
}

type AntiAffinityTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *AntiAffinityTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *AntiAffinityTestSuite) TestAntiAffinity() {
	addCaps := []v1.Capability{"FORBIDDEN_CAP", "ALLOWED_CAP"}
	dropCaps := []v1.Capability{"DROPPED_CAP"}

	s.ctx.AddMockPod(podName)
	s.addPodAndAddContainerWithCaps(addCaps, dropCaps)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				ForbiddenCapabilities: []string{"FORBIDDEN_CAP", "DROPPED_CAP"},
				Exceptions:            nil,
			},
			Diagnostics: []diagnostic.Diagnostic{
				{Message: fmt.Sprintf(addListDiagMsgFmt, containerName, "FORBIDDEN_CAP")},
				{Message: fmt.Sprintf(dropListDiagMsgFmt, containerName, dropCaps, "FORBIDDEN_CAP")},
			},
			ExpectInstantiationError: false,
		},
	})
}
