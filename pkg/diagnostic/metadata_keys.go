package diagnostic

// Metadata keys for optional diagnostic annotations.
// These are exported for use by checks and formatters.
const (
	// MetaKeyRuleID is the unique identifier for the check rule (e.g., "no-readiness-probe").
	MetaKeyRuleID = "rule_id"

	// MetaKeyFingerprint is a stable identifier for this diagnostic instance,
	// used by SARIF consumers to track issues across runs.
	MetaKeyFingerprint = "fingerprint"

	// MetaKeyRemediation is an optional remediation message specific to this diagnostic.
	// If not set, formatters will fall back to the check's default remediation.
	MetaKeyRemediation = "remediation"

	// MetaKeyCWE is the Common Weakness Enumeration identifier for security-related diagnostics.
	MetaKeyCWE = "cwe"
)
