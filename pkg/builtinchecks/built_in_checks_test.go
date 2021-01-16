package builtinchecks

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuiltInChecksWellFormed(t *testing.T) {
	checks, err := List()
	require.NoError(t, err)
	for _, check := range checks {
		t.Run(check.Name, func(t *testing.T) {
			assert.NotEmpty(t, check.Remediation, "Please add remediation")
			assert.True(t, strings.HasSuffix(check.Remediation, "."), "Please end your remediation texts with a period (got %q)", check.Remediation)
		})
	}
}
