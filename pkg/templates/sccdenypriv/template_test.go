package sccdenypriv

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/sccdenypriv/internal/params"
)

func TestSCCPriv(t *testing.T) {
	suite.Run(t, new(SCCPrivTestSuite))
}

type SCCPrivTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *SCCPrivTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *SCCPrivTestSuite) addSCCWithPriv(name string, allow_flag bool, version string) {
	s.ctx.AddMockSecurityContextConstraints(s.T(), name, allow_flag)
}

func (s *SCCPrivTestSuite) TestPrivFalse() {
	const (
		acceptableScc = "scc-priv-false"
	)

	s.addSCCWithPriv(acceptableScc, false, "v1")

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				AllowPrivilegedContainer: true,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				acceptableScc: nil,
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *SCCPrivTestSuite) TestPrivTrue() {
	const (
		unacceptableScc = "scc-priv-true"
	)

	s.addSCCWithPriv(unacceptableScc, true, "v1")

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				AllowPrivilegedContainer: true,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				unacceptableScc: {
					{Message: "SCC has allowPrivilegedContainer set to true"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}
