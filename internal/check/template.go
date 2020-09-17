package check

import (
	"golang.stackrox.io/kube-linter/internal/diagnostic"
	"golang.stackrox.io/kube-linter/internal/lintcontext"
)

// A Func is a specific lint-check, which runs on a specific objects, and emits diagnostics if problems are found.
// Checks have access to the entire LintContext, with all the objects in it, but must only report problems for the
// object passed in the second argument.
type Func func(lintCtx *lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic

// A ParameterDesc describes a parameter to a check template.
type ParameterDesc struct {
	ParamName   string
	Required    bool
	Description string
}

// ObjectKindsDesc describes a list of supported object kinds for a check template.
type ObjectKindsDesc struct {
	ObjectKinds []string `json:"objectKinds"`
}

// A Template is a template for a check.
type Template struct {
	Name                 string
	Description          string
	SupportedObjectKinds ObjectKindsDesc
	Parameters           []ParameterDesc
	Instantiate          func(params map[string]string) (Func, error)
}
