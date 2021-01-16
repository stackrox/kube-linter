package defaultchecks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.stackrox.io/kube-linter/internal/set"
	"golang.stackrox.io/kube-linter/pkg/builtinchecks"
)

func TestListReferencesOnlyValidChecks(t *testing.T) {
	allChecks, err := builtinchecks.List()
	require.NoError(t, err)
	allCheckNames := set.NewStringSet()
	for _, check := range allChecks {
		allCheckNames.Add(check.Name)
	}
	for _, defaultCheck := range List.AsSlice() {
		assert.True(t, allCheckNames.Contains(defaultCheck), "default check %s invalid", defaultCheck)
	}
}
