package diagnostic

import (
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
)

// A Diagnostic represents one specific problem diagnosed by a check.
type Diagnostic struct {
	Message string

	// TODO: add line number/col number
}

// WithContext puts a diagnostic in the context of which check emitted it,
// and which object it applied to.
type WithContext struct {
	Diagnostic  Diagnostic
	Check       string
	Remediation string
	Object      lintcontext.Object
}
