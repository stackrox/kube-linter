package params

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ugiordan/kube-chainsaw/pkg/analyzer"
	"golang.stackrox.io/kube-linter/pkg/check"
)

func TestValidateCustom_ValidParams(t *testing.T) {
	p := Params{
		Rules:        []string{},
		ExcludeRules: []string{},
		MinSeverity:  "",
	}
	err := p.ValidateCustom()
	assert.NoError(t, err)
}

func TestValidateCustom_UnknownRuleID(t *testing.T) {
	p := Params{
		Rules: []string{"KC-999", "KC-001"},
	}
	err := p.ValidateCustom()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown rule ID \"KC-999\"")
}

func TestValidateCustom_UnknownExcludeRuleID(t *testing.T) {
	p := Params{
		ExcludeRules: []string{"INVALID-RULE"},
	}
	err := p.ValidateCustom()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown excludeRules ID \"INVALID-RULE\"")
}

func TestValidateCustom_RuleInBothLists(t *testing.T) {
	known := analyzer.KnownRuleIDs()
	require.NotEmpty(t, known, "analyzer must have at least one known rule for this test")

	ruleID := known[0]
	p := Params{
		Rules:        []string{ruleID},
		ExcludeRules: []string{ruleID},
	}
	err := p.ValidateCustom()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "appears in both rules and excludeRules")
}

func TestValidateCustom_ValidMinSeverityValues(t *testing.T) {
	validSeverities := []string{"info", "note", "warning", "high", "error", "critical"}
	for _, severity := range validSeverities {
		t.Run(severity, func(t *testing.T) {
			p := Params{MinSeverity: severity}
			err := p.ValidateCustom()
			assert.NoError(t, err, "severity %q should be valid", severity)
		})
	}
}

func TestValidateCustom_InvalidMinSeverity(t *testing.T) {
	p := Params{MinSeverity: "invalid"}
	err := p.ValidateCustom()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid minSeverity \"invalid\"")
}

func TestValidateCustom_EmptyParams(t *testing.T) {
	p := Params{}
	err := p.ValidateCustom()
	assert.NoError(t, err)
}

func TestParseSeverityLevel_AllLevels(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"info", 0},
		{"note", 0},
		{"warning", 1},
		{"high", 2},
		{"error", 2},
		{"critical", 3},
		{"INFO", 0},
		{"WARNING", 1},
		{"CRITICAL", 3},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			level, err := ParseSeverityLevel(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, level)
		})
	}
}

func TestParseSeverityLevel_InvalidInput(t *testing.T) {
	_, err := ParseSeverityLevel("unknown")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid minSeverity \"unknown\"")
}

func TestParseAndValidate_ValidMap(t *testing.T) {
	m := map[string]interface{}{
		"rules":        []string{},
		"minSeverity":  "warning",
		"excludeRules": []string{},
	}

	result, err := ParseAndValidate(m)
	require.NoError(t, err)

	p, ok := result.(Params)
	require.True(t, ok)
	assert.Equal(t, "warning", p.MinSeverity)
}

func TestParseAndValidate_InvalidMap(t *testing.T) {
	m := map[string]interface{}{
		"rules":       123,
		"minSeverity": "warning",
	}

	_, err := ParseAndValidate(m)
	require.Error(t, err)
}

func TestWrapInstantiateFunc(t *testing.T) {
	called := false
	wrapped := WrapInstantiateFunc(func(p Params) (check.Func, error) {
		called = true
		assert.Equal(t, "info", p.MinSeverity)
		return nil, nil
	})

	p := Params{MinSeverity: "info"}
	_, err := wrapped(p)
	require.NoError(t, err)
	assert.True(t, called)
}
