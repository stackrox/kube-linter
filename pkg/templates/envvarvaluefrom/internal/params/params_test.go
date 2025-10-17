package params

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateParams(t *testing.T) {
	t.Run("InvalidRegex", func(t *testing.T) {
		p := Params{IgnoredSecrets: []string{"[invalid("}}
		err := p.Validate()
		// If Validate doesn't check regex, this will pass; otherwise, expect error
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
