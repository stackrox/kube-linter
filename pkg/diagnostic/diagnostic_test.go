package diagnostic

import (
	"encoding/json"
	"testing"
)

// TestDiagnostic_BackwardCompatibility verifies that the Message field
// serializes as "Message" (not "message") to maintain backward compatibility
// with existing JSON consumers.
func TestDiagnostic_BackwardCompatibility(t *testing.T) {
	d := Diagnostic{
		Message: "test message",
	}

	data, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("failed to marshal diagnostic: %v", err)
	}

	// Verify Message is capitalized (no json tag on that field)
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if _, ok := raw["Message"]; !ok {
		t.Errorf("expected JSON key 'Message' (capitalized), got keys: %v", raw)
	}
	if _, ok := raw["message"]; ok {
		t.Errorf("unexpected JSON key 'message' (lowercase) found, breaks backward compatibility")
	}
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
				Severity: "critical",
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
				Severity: "high",
				Metadata: map[string]string{
					MetaKeyFingerprint: "abc123",
				},
			},
			expectKeys:  []string{"Message", "Severity", "Metadata"},
			excludeKeys: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.diagnostic)
			if err != nil {
				t.Fatalf("failed to marshal: %v", err)
			}

			var raw map[string]interface{}
			if err := json.Unmarshal(data, &raw); err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}

			for _, key := range tt.expectKeys {
				if _, ok := raw[key]; !ok {
					t.Errorf("expected key %q in JSON, got keys: %v", key, raw)
				}
			}

			for _, key := range tt.excludeKeys {
				if _, ok := raw[key]; ok {
					t.Errorf("expected key %q to be omitted (omitempty), got keys: %v", key, raw)
				}
			}
		})
	}
}

// TestDiagnostic_MetadataKeys verifies that metadata keys are exported constants.
func TestDiagnostic_MetadataKeys(t *testing.T) {
	// Just verify the constants are accessible
	keys := []string{
		MetaKeyRuleID,
		MetaKeyFingerprint,
		MetaKeyRemediation,
		MetaKeyCWE,
	}

	for _, key := range keys {
		if key == "" {
			t.Errorf("metadata key should not be empty")
		}
	}
}
