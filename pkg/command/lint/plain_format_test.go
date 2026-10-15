package lint

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/run"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// TestPlainFormat_SeverityPrefix verifies that Diagnostic.Severity appears
// as a prefix in plain output when present.
func TestPlainFormat_SeverityPrefix(t *testing.T) {
	tests := []struct {
		name           string
		severity       string
		expectedPrefix string
	}{
		{"Critical", "critical", "[CRITICAL]"},
		{"High", "high", "[HIGH]"},
		{"Warning", "warning", "[WARNING]"},
		{"Info", "info", "[INFO]"},
		{"Empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := run.Result{
				Summary: run.Summary{
					KubeLinterVersion: "test-version",
					ChecksStatus:      run.ChecksPassed,
					CheckEndTime:      time.Now().UTC(),
				},
				Reports: []diagnostic.WithContext{
					{
						Diagnostic: diagnostic.Diagnostic{
							Message:  "test message",
							Severity: tt.severity,
						},
						Check:       "test-check",
						Remediation: "test remediation",
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
			err := plainTemplate.Execute(&buf, result)
			if err != nil {
				t.Fatalf("plainTemplate.Execute failed: %v", err)
			}

			output := buf.String()
			if tt.expectedPrefix != "" {
				if !strings.Contains(output, tt.expectedPrefix) {
					t.Errorf("expected output to contain %q, got:\n%s", tt.expectedPrefix, output)
				}
			}

			// Always verify message is present
			if !strings.Contains(output, "test message") {
				t.Errorf("expected output to contain message, got:\n%s", output)
			}
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
						Remediation: tt.checkRemediation,
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
			err := plainTemplate.Execute(&buf, result)
			if err != nil {
				t.Fatalf("plainTemplate.Execute failed: %v", err)
			}

			output := buf.String()
			if !strings.Contains(output, tt.expectedInOutput) {
				t.Errorf("expected output to contain %q, got:\n%s", tt.expectedInOutput, output)
			}
		})
	}
}
