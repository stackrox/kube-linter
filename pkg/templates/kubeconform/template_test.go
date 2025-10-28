package kubeconform

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/templates/kubeconform/internal/params"
)

func TestSliceToMap(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected map[string]struct{}
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: map[string]struct{}{},
		},
		{
			name:  "single element",
			input: []string{"pod"},
			expected: map[string]struct{}{
				"pod": {},
			},
		},
		{
			name:  "duplicate element",
			input: []string{"pod", "pod"},
			expected: map[string]struct{}{
				"pod": {},
			},
		},
		{
			name:  "multiple elements",
			input: []string{"pod", "service", "deployment"},
			expected: map[string]struct{}{
				"pod":        {},
				"service":    {},
				"deployment": {},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sliceToMap(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

//go:embed testdata/demonset.yaml
var invalidDaemonSetYAML string

func TestValidateWithEmptyParams(t *testing.T) {
	object := lintcontext.Object{
		Metadata: lintcontext.ObjectMetadata{
			FilePath: "demonset.yaml",
			Raw:      []byte(invalidDaemonSetYAML),
		},
	}

	// Test with empty params (default kubeconform settings)
	emptyParams := params.Params{
		SchemaLocations: []string{"testdata/schema.json"},
		Strict:          true,
	}

	checkFunc, err := validate(emptyParams)
	assert.NoError(t, err)
	assert.NotNil(t, checkFunc)

	// Run the validation (LintContext not used by kubeconform validator)
	var ctx lintcontext.LintContext
	diagnostics := checkFunc(ctx, object)

	require.Len(t, diagnostics, 1)
	assert.Contains(t, diagnostics[0].Message, `- at '/spec/selector/matchLabels/k8s-app': got boolean, want null or string - at '/spec': additional properties 'replicas' not allowed`)
}
