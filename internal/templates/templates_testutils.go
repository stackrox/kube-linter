package templates

import (
	"sort"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/diagnostic"
	"golang.stackrox.io/kube-linter/internal/lintcontext"
)

// TemplateTestSuite is a basic testing suite for all templates
// test with some generic helper functions
type TemplateTestSuite struct {
	suite.Suite

	Template check.Template
}

// TestCase represents a single test case which can be verified under a LintContext
type TestCase struct {
	Param       interface{}
	Diagnostics []diagnostic.Diagnostic
	ExpectError bool
}

// Init initializes the test suite with a template
func (s *TemplateTestSuite) Init(templateKey string) {
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
		checkFunc, err := s.Template.Instantiate(c.Param)
		if c.ExpectError {
			s.Error(err, "param should have caused error but did not raise one")
			continue
		}
		for _, obj := range ctx.GetObjects() {
			diagnostics := checkFunc(ctx, obj)
			passed := s.compareDiagnostics(c.Diagnostics, diagnostics)
			s.True(passed, "expected diagnostics: %q. actual: %q", c.Diagnostics, diagnostics)
		}
	}
}

func (s *TemplateTestSuite) compareDiagnostics(expected, actual []diagnostic.Diagnostic) bool {
	if len(expected) != len(actual) {
		return false
	}
	sort.Slice(expected, func(i, j int) bool {
		return expected[i].Message < expected[j].Message
	})
	sort.Slice(actual, func(i, j int) bool {
		return actual[i].Message < actual[j].Message
	})

	for i := 0; i < len(expected); i++ {
		if expected[i].Message != actual[i].Message {
			return false
		}
	}

	return true
}
