package lint

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			name:        "multiple formats no outputs (rejected)",
			formats:     []string{"json", "sarif"},
			outputs:     []string{},
			wantErr:     true,
			errContains: "multiple formats require explicit --output flags",
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

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			require.NoError(t, err)
			assert.Len(t, pairs, tt.wantPairCount)

			// Verify pairing
			for i, pair := range pairs {
				assert.Equal(t, tt.formats[i], string(pair.Format))

				expectedOutput := ""
				if len(tt.outputs) > 0 {
					expectedOutput = tt.outputs[i]
				}
				assert.Equal(t, expectedOutput, pair.Output)

				if tt.wantStdout {
					assert.Empty(t, pair.Output, "Pair %d: expected stdout (empty output)", i)
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

	require.NoError(t, err)
	assert.Equal(t, string(common.JSONFormat), string(pairs[0].Format))
	assert.Equal(t, string(common.SARIFFormat), string(pairs[1].Format))
}
