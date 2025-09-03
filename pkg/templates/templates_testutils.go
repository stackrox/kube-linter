package templates

import (
	"fmt"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
)

// TemplateTestSuite is a basic testing suite for all templates
// test with some generic helper functions
type TemplateTestSuite struct {
	suite.Suite

	Template check.Template
}

// TestCase represents a single test case which can be verified under a LintContext
type TestCase struct {
	Param                    interface{}
	Diagnostics              map[string][]diagnostic.Diagnostic
	ExpectInstantiationError bool
}

// Init initializes the test suite with a template
func (s *TemplateTestSuite) Init(templateKey string) {
	s.T().Helper()
	t, ok := Get(templateKey)
	s.True(ok, "template with key %q not found", templateKey)
	s.Template = t
}

// Validate validates the given test cases against the LintContext passed in.
func (s *TemplateTestSuite) Validate(
	ctx lintcontext.LintContext,
	cases []TestCase,
) {
	for _, c := range cases {
		s.Run(fmt.Sprintf("%+v", c.Param), func() {
			checkFunc, err := s.Template.Instantiate(c.Param)
			if c.ExpectInstantiationError {
				s.Error(err, "param should have caused error but did not raise one")
				return
			}
			s.Require().NoError(err)
			for _, obj := range ctx.Objects() {
				diagnostics := checkFunc(ctx, obj)
				s.compareDiagnostics(c.Diagnostics[obj.K8sObject.GetName()], diagnostics)
			}
		})
	}
}

func (s *TemplateTestSuite) compareDiagnostics(expected, actual []diagnostic.Diagnostic) {
	s.T().Helper()
	expectedMessages, actualMessages := make([]string, 0, len(expected)), make([]string, 0, len(actual))
	for _, diag := range expected {
		expectedMessages = append(expectedMessages, diag.Message)
	}
	for _, diag := range actual {
		actualMessages = append(actualMessages, diag.Message)
	}
	s.ElementsMatch(expectedMessages, actualMessages, "expected diagnostics and actual diagnostics do not match")
}
