package diagnostic

import (
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
)

// A Diagnostic represents one specific problem diagnosed by a check.
type Diagnostic struct {
	Message string

	// Severity is the severity level for this diagnostic (e.g., "critical", "high", "warning", "info").
	// Optional. If not set, formatters may use a default severity.
	Severity string `json:",omitempty"`

	// Metadata is an optional key-value map for additional diagnostic information.
	// Common keys are defined as MetaKey* constants in metadata_keys.go.
	Metadata map[string]string `json:",omitempty"`

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
