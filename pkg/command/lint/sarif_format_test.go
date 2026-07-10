package lint

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/run"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// TestFormatSarif_SeverityMapping verifies that Diagnostic.Severity is mapped
// to SARIF result.level correctly.
func TestFormatSarif_SeverityMapping(t *testing.T) {
	tests := []struct {
		name             string
		diagnosticSev    string
		expectedSARIFLvl string
	}{
		{"Critical", "critical", "error"},
		{"High", "high", "error"},
		{"Warning", "warning", "warning"},
		{"Info", "info", "note"},
		{"Empty", "", "warning"}, // default
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := run.Result{
				Summary: run.Summary{
					KubeLinterVersion: "test-version",
					ChecksStatus:      run.ChecksPassed,
					CheckEndTime:      time.Now().UTC(),
				},
				Checks: []config.Check{
					{
						Name:        "test-check",
						Description: "test description",
						Remediation: "test remediation",
						Template:    "readiness-probe",
					},
				},
				Reports: []diagnostic.WithContext{
					{
						Diagnostic: diagnostic.Diagnostic{
							Message:  "test message",
							Severity: tt.diagnosticSev,
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
			err := formatSarif(&buf, result)
			if err != nil {
				t.Fatalf("formatSarif failed: %v", err)
			}

			var sarifOutput map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &sarifOutput); err != nil {
				t.Fatalf("failed to unmarshal SARIF: %v", err)
			}

			// Navigate to results[0].level
			runs := sarifOutput["runs"].([]interface{})
			run := runs[0].(map[string]interface{})
			results := run["results"].([]interface{})
			result0 := results[0].(map[string]interface{})

			actualLevel, ok := result0["level"].(string)
			if !ok {
				t.Fatalf("expected level field in SARIF result, got: %v", result0)
			}
			if actualLevel != tt.expectedSARIFLvl {
				t.Errorf("expected SARIF level %q, got %q", tt.expectedSARIFLvl, actualLevel)
			}
		})
	}
}

// TestFormatSarif_MetadataInProperties verifies that Diagnostic.Metadata
// is serialized into SARIF result.properties bag.
func TestFormatSarif_MetadataInProperties(t *testing.T) {
	result := run.Result{
		Summary: run.Summary{
			KubeLinterVersion: "test-version",
			ChecksStatus:      run.ChecksPassed,
			CheckEndTime:      time.Now().UTC(),
		},
		Checks: []config.Check{
			{
				Name:        "test-check",
				Description: "test description",
				Remediation: "test remediation",
				Template:    "readiness-probe",
			},
		},
		Reports: []diagnostic.WithContext{
			{
				Diagnostic: diagnostic.Diagnostic{
					Message:  "test message",
					Severity: "high",
					Metadata: map[string]string{
						diagnostic.MetaKeyRuleID:      "no-readiness-probe",
						diagnostic.MetaKeyFingerprint: "abc123",
						diagnostic.MetaKeyRemediation: "custom remediation",
						diagnostic.MetaKeyCWE:         "CWE-778",
						"custom_key":                  "custom_value",
					},
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
	err := formatSarif(&buf, result)
	if err != nil {
		t.Fatalf("formatSarif failed: %v", err)
	}

	var sarifOutput map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &sarifOutput); err != nil {
		t.Fatalf("failed to unmarshal SARIF: %v", err)
	}

	// Navigate to results[0].properties
	runs := sarifOutput["runs"].([]interface{})
	run := runs[0].(map[string]interface{})
	results := run["results"].([]interface{})
	result0 := results[0].(map[string]interface{})

	properties, ok := result0["properties"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected properties bag in SARIF result")
	}

	// Verify metadata is present
	expectedKeys := []string{"rule_id", "fingerprint", "remediation", "cwe", "custom_key"}
	for _, key := range expectedKeys {
		if _, ok := properties[key]; !ok {
			t.Errorf("expected properties[%q] to be present, got properties: %v", key, properties)
		}
	}

	if properties["rule_id"] != "no-readiness-probe" {
		t.Errorf("expected rule_id=no-readiness-probe, got %v", properties["rule_id"])
	}
}

// TestFormatSarif_PartialFingerprints verifies that when Diagnostic.Metadata
// contains a fingerprint, it's added to SARIF result.partialFingerprints.
func TestFormatSarif_PartialFingerprints(t *testing.T) {
	result := run.Result{
		Summary: run.Summary{
			KubeLinterVersion: "test-version",
			ChecksStatus:      run.ChecksPassed,
			CheckEndTime:      time.Now().UTC(),
		},
		Checks: []config.Check{
			{
				Name:        "test-check",
				Description: "test description",
				Remediation: "test remediation",
				Template:    "readiness-probe",
			},
		},
		Reports: []diagnostic.WithContext{
			{
				Diagnostic: diagnostic.Diagnostic{
					Message:  "test message",
					Severity: "critical",
					Metadata: map[string]string{
						diagnostic.MetaKeyFingerprint: "stable-fp-123",
					},
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
	err := formatSarif(&buf, result)
	if err != nil {
		t.Fatalf("formatSarif failed: %v", err)
	}

	var sarifOutput map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &sarifOutput); err != nil {
		t.Fatalf("failed to unmarshal SARIF: %v", err)
	}

	// Navigate to results[0].partialFingerprints
	runs := sarifOutput["runs"].([]interface{})
	run := runs[0].(map[string]interface{})
	results := run["results"].([]interface{})
	result0 := results[0].(map[string]interface{})

	partialFingerprints, ok := result0["partialFingerprints"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected partialFingerprints in SARIF result")
	}

	if partialFingerprints["primaryLocationLineHash"] != "stable-fp-123" {
		t.Errorf("expected primaryLocationLineHash=stable-fp-123, got %v", partialFingerprints["primaryLocationLineHash"])
	}
}

// TestFormatSarif_NoMetadata verifies that when Diagnostic has no metadata,
// no properties/partialFingerprints are added.
func TestFormatSarif_NoMetadata(t *testing.T) {
	result := run.Result{
		Summary: run.Summary{
			KubeLinterVersion: "test-version",
			ChecksStatus:      run.ChecksPassed,
			CheckEndTime:      time.Now().UTC(),
		},
		Checks: []config.Check{
			{
				Name:        "test-check",
				Description: "test description",
				Remediation: "test remediation",
				Template:    "readiness-probe",
			},
		},
		Reports: []diagnostic.WithContext{
			{
				Diagnostic: diagnostic.Diagnostic{
					Message: "test message",
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
	err := formatSarif(&buf, result)
	if err != nil {
		t.Fatalf("formatSarif failed: %v", err)
	}

	var sarifOutput map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &sarifOutput); err != nil {
		t.Fatalf("failed to unmarshal SARIF: %v", err)
	}

	// Navigate to results[0]
	runs := sarifOutput["runs"].([]interface{})
	run := runs[0].(map[string]interface{})
	results := run["results"].([]interface{})
	result0 := results[0].(map[string]interface{})

	// properties and partialFingerprints should not be present
	if _, ok := result0["properties"]; ok {
		t.Errorf("expected no properties when Metadata is nil")
	}
	if _, ok := result0["partialFingerprints"]; ok {
		t.Errorf("expected no partialFingerprints when Metadata has no fingerprint")
	}
}
