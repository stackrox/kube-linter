package params

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name          string
		params        Params
		expectedError error
	}{
		{
			name: "valid annotation",
			params: Params{
				Annotation: "some-annotation",
			},
			expectedError: nil,
		},
		{
			name: "missing annotation",
			params: Params{
				Annotation: "",
			},
			expectedError: errors.New("invalid parameters: required param annotation not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError.Error())
			}
		})
	}
}

func TestParseAndValidate(t *testing.T) {
	t.Run("valid map input", func(t *testing.T) {
		m := map[string]interface{}{
			"annotation": "required-annotation",
		}
		result, err := ParseAndValidate(m)
		assert.NoError(t, err)

		params, ok := result.(Params)
		assert.True(t, ok)
		assert.Equal(t, "required-annotation", params.Annotation)
	})

	t.Run("missing annotation in map", func(t *testing.T) {
		m := map[string]interface{}{}
		_, err := ParseAndValidate(m)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "required param annotation not found")
	})
}

func TestWrapInstantiateFunc(t *testing.T) {
	mockFunc := func(p Params) (check.Func, error) {
		return func(ctx lintcontext.LintContext, obj lintcontext.Object) []diagnostic.Diagnostic {
			return []diagnostic.Diagnostic{
				{
					Message: "mocked diagnostic",
				},
			}
		}, nil
	}

	wrapped := WrapInstantiateFunc(mockFunc)

	t.Run("wrapped function works", func(t *testing.T) {
		params := Params{Annotation: "test-annotation"}
		fn, err := wrapped(params)
		assert.NoError(t, err)

		var dummyCtx lintcontext.LintContext
		var dummyObj lintcontext.Object

		diagnostics := fn(dummyCtx, dummyObj)
		assert.Len(t, diagnostics, 1)
		assert.Equal(t, "mocked diagnostic", diagnostics[0].Message)
	})
}
