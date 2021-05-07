package lintcontext

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateContextsFromHelmArchive(t *testing.T) {
	fileName := "../../tests/testdata/mychart-0.1.0.tgz"
	file, err := os.Open(fileName)
	require.NoError(t, err)

	lintCtxs, err := CreateContextsFromHelmArchive("test", file)
	assert.NoError(t, err)

	var atLeastOneObjectFound bool
	for _, lintCtx := range lintCtxs {
		if len(lintCtx.Objects()) > 0 {
			atLeastOneObjectFound = true
			break
		}
	}
	assert.True(t, atLeastOneObjectFound, "no valid objects found")
}
