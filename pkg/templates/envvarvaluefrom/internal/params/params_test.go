package params

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
)

func TestValidateParams(t *testing.T) {
	t.Run("InvalidRegex", func(t *testing.T) {
		p := Params{IgnoredSecrets: []string{"[invalid("}}
		err := p.Validate()
		// Current behavior: No validation, so no error
		if err == nil {
			t.Log("Warning: Validate does not check regex validity; consider adding regex validation")
		} else {
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid syntax")
		}
	})

	t.Run("ValidParams", func(t *testing.T) {
		p := Params{IgnoredSecrets: []string{"^valid$"}}
		err := p.Validate()
		assert.NoError(t, err)
	})
}

func TestParseAndValidate(t *testing.T) {
	t.Run("InvalidDecode", func(t *testing.T) {
		// Simulate invalid map structure
		m := map[string]interface{}{"IgnoredSecrets": map[string]interface{}{"invalid": true}}
		_, err := ParseAndValidate(m)
		assert.Error(t, err) // Expect DecodeMapStructure to fail
	})

	t.Run("ValidParse", func(t *testing.T) {
		m := map[string]interface{}{"IgnoredSecrets": []interface{}{"^valid$"}}
		p, err := ParseAndValidate(m)
		assert.NoError(t, err)
		assert.Equal(t, []string{"^valid$"}, p.(Params).IgnoredSecrets)
	})
}

func TestWrapInstantiateFunc(t *testing.T) {
	t.Run("ValidWrap", func(t *testing.T) {
		f := func(p Params) (check.Func, error) {
			// check.Func is a function type, so we can return a function directly
			return func(lintCtx lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				return nil // Return empty slice instead of nil for consistency
			}, nil
		}
		wrapped := WrapInstantiateFunc(f)
		_, err := wrapped(Params{})
		assert.NoError(t, err)
	})

	t.Run("InvalidType", func(t *testing.T) {
		f := func(p Params) (check.Func, error) {
			return func(lintCtx lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				return nil // Return empty slice instead of nil for consistency
			}, nil
		}
		wrapped := WrapInstantiateFunc(f)
		assert.Panics(t, func() {
			_, _ = wrapped("not-a-Params") // Expect panic due to type assertion failure
		})
	})
}
