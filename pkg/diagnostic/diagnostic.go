package diagnostic

import (
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
)

// Severity is the severity level of a diagnostic. Custom values remain
// possible, but producers should use the Severity* constants.
type Severity string

const (
	// SeverityCritical is the highest severity level.
	SeverityCritical Severity = "critical"
	// SeverityHigh is used for severe findings that are not critical.
	SeverityHigh Severity = "high"
	// SeverityWarning is used for findings that should be reviewed.
	SeverityWarning Severity = "warning"
	// SeverityInfo is used for informational findings.
	SeverityInfo Severity = "info"
)

// A Diagnostic represents one specific problem diagnosed by a check.
type Diagnostic struct {
	Message string

	// Severity is the severity level for this diagnostic.
	// Optional. If not set, formatters may use a default severity.
	Severity Severity `json:",omitempty"`

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
