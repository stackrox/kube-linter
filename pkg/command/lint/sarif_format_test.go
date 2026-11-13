package lint

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/run"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func testRunResult(diag diagnostic.Diagnostic) run.Result {
	return run.Result{
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
				Diagnostic:  diag,
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
}

// firstSarResult formats result into SARIF and returns runs[0].results[0].
func firstSarResult(t *testing.T, result run.Result) map[string]interface{} {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, formatSarif(&buf, result))

	var sarifOutput map[string]interface{}
	require.NoError(t, json.Unmarshal(buf.Bytes(), &sarifOutput))

	runs, ok := sarifOutput["runs"].([]interface{})
	require.True(t, ok, "expected runs in SARIF output")
	sarifRun, ok := runs[0].(map[string]interface{})
	require.True(t, ok, "expected run object in SARIF output")
	results, ok := sarifRun["results"].([]interface{})
	require.True(t, ok, "expected results in SARIF output")
	result0, ok := results[0].(map[string]interface{})
	require.True(t, ok, "expected result object in SARIF output")
	return result0
}

// TestFormatSarif_SeverityMapping verifies that Diagnostic.Severity is mapped
// to SARIF result.level correctly.
func TestFormatSarif_SeverityMapping(t *testing.T) {
	tests := []struct {
		name             string
		diagnosticSev    diagnostic.Severity
		expectedSARIFLvl string
	}{
		{"Critical", diagnostic.SeverityCritical, "error"},
		{"High", diagnostic.SeverityHigh, "error"},
		{"Warning", diagnostic.SeverityWarning, "warning"},
		{"Info", diagnostic.SeverityInfo, "note"},
		{"Empty", "", "warning"}, // default
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result0 := firstSarResult(t, testRunResult(diagnostic.Diagnostic{
				Message:  "test message",
				Severity: tt.diagnosticSev,
			}))

			assert.Equal(t, tt.expectedSARIFLvl, result0["level"])
		})
	}
}

// TestMapSeverityToSARIFLevel verifies the severity-to-SARIF level mapping,
// including the fallback for unrecognized severity values.
func TestMapSeverityToSARIFLevel(t *testing.T) {
	tests := []struct {
		severity diagnostic.Severity
		expected string
	}{
		{diagnostic.SeverityCritical, "error"},
		{diagnostic.SeverityHigh, "error"},
		{diagnostic.SeverityWarning, "warning"},
		{diagnostic.SeverityInfo, "note"},
		{"CRITICAL", "error"},
		{"unknown-severity", "warning"},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.expected, mapSeverityToSARIFLevel(tt.severity), "severity %q", tt.severity)
	}
}

// TestFormatSarif_MetadataInProperties verifies that Diagnostic.Metadata
// is serialized into SARIF result.properties bag.
func TestFormatSarif_MetadataInProperties(t *testing.T) {
	result0 := firstSarResult(t, testRunResult(diagnostic.Diagnostic{
		Message:  "test message",
		Severity: diagnostic.SeverityHigh,
		Metadata: map[string]string{
			diagnostic.MetaKeyRuleID:      "no-readiness-probe",
			diagnostic.MetaKeyFingerprint: "abc123",
			diagnostic.MetaKeyRemediation: "custom remediation",
			diagnostic.MetaKeyCWE:         "CWE-778",
			"custom_key":                  "custom_value",
		},
	}))

	properties, ok := result0["properties"].(map[string]interface{})
	require.True(t, ok, "expected properties bag in SARIF result")

	for _, key := range []string{"rule_id", "fingerprint", "remediation", "cwe", "custom_key"} {
		assert.Contains(t, properties, key)
	}
	assert.Equal(t, "no-readiness-probe", properties["rule_id"])
}

// TestFormatSarif_PartialFingerprints verifies that when Diagnostic.Metadata
// contains a fingerprint, it's added to SARIF result.partialFingerprints.
func TestFormatSarif_PartialFingerprints(t *testing.T) {
	result0 := firstSarResult(t, testRunResult(diagnostic.Diagnostic{
		Message:  "test message",
		Severity: diagnostic.SeverityCritical,
		Metadata: map[string]string{
			diagnostic.MetaKeyFingerprint: "stable-fp-123",
		},
	}))

	partialFingerprints, ok := result0["partialFingerprints"].(map[string]interface{})
	require.True(t, ok, "expected partialFingerprints in SARIF result")
	assert.Equal(t, "stable-fp-123", partialFingerprints["primaryLocationLineHash"])
}

// TestFormatSarif_NoMetadata verifies that when Diagnostic has no metadata,
// no properties/partialFingerprints are added.
func TestFormatSarif_NoMetadata(t *testing.T) {
	result0 := firstSarResult(t, testRunResult(diagnostic.Diagnostic{
		Message: "test message",
	}))

	assert.NotContains(t, result0, "properties")
	assert.NotContains(t, result0, "partialFingerprints")
}

// TestFormatSarif_TagsMetadataAsArray verifies that a metadata key named
// "tags" is serialized as an array, because the SARIF property bag reserves
// "tags" as an array of strings.
func TestFormatSarif_TagsMetadataAsArray(t *testing.T) {
	result0 := firstSarResult(t, testRunResult(diagnostic.Diagnostic{
		Message:  "test message",
		Severity: diagnostic.SeverityHigh,
		Metadata: map[string]string{
			"tags":                   "security",
			diagnostic.MetaKeyRuleID: "no-readiness-probe",
			diagnostic.MetaKeyCWE:    "CWE-778",
		},
	}))

	properties, ok := result0["properties"].(map[string]interface{})
	require.True(t, ok, "expected properties bag in SARIF result")

	tags, ok := properties["tags"].([]interface{})
	require.True(t, ok, "expected properties.tags to be an array, got %T", properties["tags"])
	assert.Equal(t, []interface{}{"security"}, tags)

	assert.Equal(t, "no-readiness-probe", properties[diagnostic.MetaKeyRuleID])
}

// TestFormatLintSarif verifies the format wrapper forwards run.Result data.
func TestFormatLintSarif(t *testing.T) {
	var buf bytes.Buffer
	require.NoError(t, formatLintSarif(&buf, testRunResult(diagnostic.Diagnostic{Message: "test message"})))
	assert.NotEmpty(t, buf.Bytes())
}

// TestFormatLintSarif_WrongDataType verifies the wrapper rejects data that is
// not a run.Result.
func TestFormatLintSarif_WrongDataType(t *testing.T) {
	assert.Error(t, formatLintSarif(&bytes.Buffer{}, "not a result"))
}

// TestFormatSarif_UnknownTemplate verifies that a check referencing a
// non-existent template fails rule generation instead of silently producing
// an incomplete SARIF report.
func TestFormatSarif_UnknownTemplate(t *testing.T) {
	result := run.Result{
		Summary: run.Summary{
			KubeLinterVersion: "test-version",
			ChecksStatus:      run.ChecksPassed,
			CheckEndTime:      time.Now().UTC(),
		},
		Checks: []config.Check{
			{
				Name:        "broken-check",
				Description: "test description",
				Remediation: "test remediation",
				Template:    "no-such-template",
			},
		},
	}

	err := formatSarif(&bytes.Buffer{}, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "non-existent template")
}
