package lint

import (
	"testing"

	"golang.stackrox.io/kube-linter/pkg/command/common"
)

func TestValidateAndPairFormatsOutputs(t *testing.T) {
	allowedFormats := []string{"json", "sarif", "plain", "markdown"}

	tests := []struct {
		name          string
		formats       []string
		outputs       []string
		wantErr       bool
		wantPairCount int
		wantStdout    bool
		errContains   string
	}{
		{
			name:          "single format no output (stdout)",
			formats:       []string{"json"},
			outputs:       []string{},
			wantErr:       false,
			wantPairCount: 1,
			wantStdout:    true,
		},
		{
			name:          "multiple formats no outputs (all stdout)",
			formats:       []string{"json", "sarif"},
			outputs:       []string{},
			wantErr:       false,
			wantPairCount: 2,
			wantStdout:    true,
		},
		{
			name:          "formats and outputs matching",
			formats:       []string{"json", "sarif"},
			outputs:       []string{"out.json", "out.sarif"},
			wantErr:       false,
			wantPairCount: 2,
			wantStdout:    false,
		},
		{
			name:        "format/output count mismatch",
			formats:     []string{"json", "sarif"},
			outputs:     []string{"out.json"},
			wantErr:     true,
			errContains: "format/output mismatch",
		},
		{
			name:        "duplicate output files",
			formats:     []string{"json", "sarif"},
			outputs:     []string{"out.txt", "out.txt"},
			wantErr:     true,
			errContains: "duplicate output file",
		},
		{
			name:        "invalid format",
			formats:     []string{"json", "invalid"},
			outputs:     []string{},
			wantErr:     true,
			errContains: "invalid format",
		},
		{
			name:        "no formats specified",
			formats:     []string{},
			outputs:     []string{},
			wantErr:     true,
			errContains: "at least one format must be specified",
		},
		{
			name:          "three formats three outputs",
			formats:       []string{"json", "sarif", "plain"},
			outputs:       []string{"out1.json", "out2.sarif", "out3.txt"},
			wantErr:       false,
			wantPairCount: 3,
			wantStdout:    false,
		},
		{
			name:          "duplicate formats allowed",
			formats:       []string{"json", "json"},
			outputs:       []string{"out1.json", "out2.json"},
			wantErr:       false,
			wantPairCount: 2,
			wantStdout:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pairs, err := ValidateAndPairFormatsOutputs(tt.formats, tt.outputs, allowedFormats)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAndPairFormatsOutputs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err != nil && tt.errContains != "" {
					if !contains(err.Error(), tt.errContains) {
						t.Errorf("Expected error to contain %q, got: %v", tt.errContains, err)
					}
				}
				return
			}

			if len(pairs) != tt.wantPairCount {
				t.Errorf("Expected %d pairs, got %d", tt.wantPairCount, len(pairs))
			}

			// Verify pairing
			for i, pair := range pairs {
				if string(pair.Format) != tt.formats[i] {
					t.Errorf("Pair %d: expected format %s, got %s", i, tt.formats[i], pair.Format)
				}

				expectedOutput := ""
				if len(tt.outputs) > 0 {
					expectedOutput = tt.outputs[i]
				}

				if pair.Output != expectedOutput {
					t.Errorf("Pair %d: expected output %s, got %s", i, expectedOutput, pair.Output)
				}

				if tt.wantStdout && pair.Output != "" {
					t.Errorf("Pair %d: expected stdout (empty output), got %s", i, pair.Output)
				}
			}
		})
	}
}

func TestValidateAndPairFormatsOutputs_FormatTypes(t *testing.T) {
	allowedFormats := []string{string(common.JSONFormat), string(common.SARIFFormat), string(common.PlainFormat)}

	pairs, err := ValidateAndPairFormatsOutputs(
		[]string{"json", "sarif"},
		[]string{"out.json", "out.sarif"},
		allowedFormats,
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if pairs[0].Format != common.JSONFormat {
		t.Errorf("Expected first format to be JSONFormat, got %v", pairs[0].Format)
	}

	if pairs[1].Format != common.SARIFFormat {
		t.Errorf("Expected second format to be SARIFFormat, got %v", pairs[1].Format)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || substr == "" || containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
