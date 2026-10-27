package params

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
)

func TestValidate(t *testing.T) {
	empty := Params{}
	assert.NoError(t, empty.Validate())

	withAllowList := Params{AllowList: []string{"^localhost:5000/"}}
	assert.NoError(t, withAllowList.Validate())
}

func TestParseAndValidate(t *testing.T) {
	t.Run("with allow list", func(t *testing.T) {
		result, err := ParseAndValidate(map[string]interface{}{
			"allowList": []interface{}{"^localhost:5000/", ".*internal.*"},
		})
		require.NoError(t, err)

		params, ok := result.(Params)
		require.True(t, ok)
		assert.Equal(t, []string{"^localhost:5000/", ".*internal.*"}, params.AllowList)
	})

	t.Run("empty map", func(t *testing.T) {
		result, err := ParseAndValidate(map[string]interface{}{})
		require.NoError(t, err)

		params, ok := result.(Params)
		require.True(t, ok)
		assert.Empty(t, params.AllowList)
	})

	t.Run("wrong type for allow list", func(t *testing.T) {
		_, err := ParseAndValidate(map[string]interface{}{
			"allowList": 42,
		})
		assert.Error(t, err)
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
		params := Params{AllowList: []string{"^localhost:5000/"}}
		fn, err := wrapped(params)
		assert.NoError(t, err)

		var dummyCtx lintcontext.LintContext
		var dummyObj lintcontext.Object

		diagnostics := fn(dummyCtx, dummyObj)
		assert.Len(t, diagnostics, 1)
		assert.Equal(t, "mocked diagnostic", diagnostics[0].Message)
	})
}
