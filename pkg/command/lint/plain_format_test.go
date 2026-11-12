package lint

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/run"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func executePlainTemplate(t *testing.T, diag diagnostic.Diagnostic, remediation string) string {
	t.Helper()
	result := run.Result{
		Summary: run.Summary{
			KubeLinterVersion: "test-version",
			ChecksStatus:      run.ChecksPassed,
			CheckEndTime:      time.Now().UTC(),
		},
		Reports: []diagnostic.WithContext{
			{
				Diagnostic:  diag,
				Check:       "test-check",
				Remediation: remediation,
				Object: lintcontext.Object{
					Metadata: lintcontext.ObjectMetadata{
						FilePath: "test.yaml",
					},
					K8sObject: &unstructured.Unstructured{
						Object: map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "Pod",
							"metadata": map[string]interface{}{
								"name":      "test-pod",
								"namespace": "default",
							},
						},
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	require.NoError(t, plainTemplate.Execute(&buf, result))
	return buf.String()
}

// TestPlainFormat_SeverityPrefix verifies that Diagnostic.Severity appears
// as a prefix in plain output when present.
func TestPlainFormat_SeverityPrefix(t *testing.T) {
	tests := []struct {
		name           string
		severity       diagnostic.Severity
		expectedPrefix string
	}{
		{"Critical", diagnostic.SeverityCritical, "[CRITICAL]"},
		{"High", diagnostic.SeverityHigh, "[HIGH]"},
		{"Warning", diagnostic.SeverityWarning, "[WARNING]"},
		{"Info", diagnostic.SeverityInfo, "[INFO]"},
		{"Empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := executePlainTemplate(t, diagnostic.Diagnostic{
				Message:  "test message",
				Severity: tt.severity,
			}, "test remediation")

			if tt.expectedPrefix != "" {
				assert.Contains(t, output, tt.expectedPrefix)
			}

			// Always verify message is present
			assert.Contains(t, output, "test message")
		})
	}
}

// TestPlainFormat_EffectiveRemediation verifies that when Diagnostic.Metadata
// contains a remediation, it overrides the check's default remediation.
func TestPlainFormat_EffectiveRemediation(t *testing.T) {
	tests := []struct {
		name                string
		checkRemediation    string
		metadataRemediation string
		expectedInOutput    string
	}{
		{
			name:                "UseCheckRemediation",
			checkRemediation:    "default remediation",
			metadataRemediation: "",
			expectedInOutput:    "default remediation",
		},
		{
			name:                "UseMetadataRemediation",
			checkRemediation:    "default remediation",
			metadataRemediation: "custom remediation",
			expectedInOutput:    "custom remediation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diag := diagnostic.Diagnostic{
				Message: "test message",
			}
			if tt.metadataRemediation != "" {
				diag.Metadata = map[string]string{
					diagnostic.MetaKeyRemediation: tt.metadataRemediation,
				}
			}

			output := executePlainTemplate(t, diag, tt.checkRemediation)
			assert.Contains(t, output, tt.expectedInOutput)
		})
	}
}
