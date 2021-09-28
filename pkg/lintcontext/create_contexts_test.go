package lintcontext

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const chartTarball = "../../tests/testdata/mychart-0.1.0.tgz"

func TestCreateContextsFromHelmArchive(t *testing.T) {
	file, err := os.Open(chartTarball)
	require.NoError(t, err)

	defer func() {
		require.NoError(t, file.Close())
	}()

	lintCtxs, err := CreateContextsFromHelmArchive("test", file)
	require.NoError(t, err)

	assert.NotNil(t, verifyAndGetContext(t, lintCtxs))
}

func TestCreateContexts_WithHelmArchive(t *testing.T) {
	lintCtxs, err := CreateContexts(chartTarball)
	require.NoError(t, err)

	assert.NotNil(t, verifyAndGetContext(t, lintCtxs))
}

func verifyAndGetContext(t *testing.T, lintCtxs []LintContext) LintContext {
	assert.Len(t, lintCtxs, 1, "expecting single lint context to be present")
	lintCtx := lintCtxs[0]

	assert.NotEmpty(t, lintCtx.Objects(), "no valid objects found")
	assert.Empty(t, lintCtx.InvalidObjects(), "no invalid objects expected")

	return lintCtx
}
