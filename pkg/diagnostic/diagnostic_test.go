package diagnostic

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDiagnostic_BackwardCompatibility verifies that the Message field
// serializes as "Message" (not "message") to maintain backward compatibility
// with existing JSON consumers.
func TestDiagnostic_BackwardCompatibility(t *testing.T) {
	d := Diagnostic{Message: "test message"}

	data, err := json.Marshal(d)
	require.NoError(t, err)

	var raw map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &raw))

	assert.Contains(t, raw, "Message")
	assert.NotContains(t, raw, "message")
}

// TestDiagnostic_SeverityAndMetadata verifies that Severity and Metadata
// fields serialize correctly and are omitted when empty.
func TestDiagnostic_SeverityAndMetadata(t *testing.T) {
	tests := []struct {
		name        string
		diagnostic  Diagnostic
		expectKeys  []string
		excludeKeys []string
	}{
		{
			name: "EmptyDiagnostic",
			diagnostic: Diagnostic{
				Message: "test",
			},
			expectKeys:  []string{"Message"},
			excludeKeys: []string{"Severity", "Metadata"},
		},
		{
			name: "WithSeverity",
			diagnostic: Diagnostic{
				Message:  "test",
				Severity: SeverityCritical,
			},
			expectKeys:  []string{"Message", "Severity"},
			excludeKeys: []string{"Metadata"},
		},
		{
			name: "WithMetadata",
			diagnostic: Diagnostic{
				Message: "test",
				Metadata: map[string]string{
					MetaKeyRuleID: "no-readiness-probe",
				},
			},
			expectKeys:  []string{"Message", "Metadata"},
			excludeKeys: []string{"Severity"},
		},
		{
			name: "WithBoth",
			diagnostic: Diagnostic{
				Message:  "test",
				Severity: SeverityHigh,
				Metadata: map[string]string{
					MetaKeyFingerprint: "abc123",
				},
			},
			expectKeys: []string{"Message", "Severity", "Metadata"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.diagnostic)
			require.NoError(t, err)

			var raw map[string]interface{}
			require.NoError(t, json.Unmarshal(data, &raw))

			for _, key := range tt.expectKeys {
				assert.Contains(t, raw, key)
			}

			for _, key := range tt.excludeKeys {
				assert.NotContains(t, raw, key)
			}
		})
	}
}

// TestDiagnostic_MetadataKeys verifies that metadata keys are exported constants.
func TestDiagnostic_MetadataKeys(t *testing.T) {
	for _, key := range []string{
		MetaKeyRuleID,
		MetaKeyFingerprint,
		MetaKeyRemediation,
		MetaKeyCWE,
	} {
		assert.NotEmpty(t, key)
	}
}
