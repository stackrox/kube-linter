package lintcontext

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const chartTarball = "../../tests/testdata/mychart-0.1.0.tgz"

func TestCreateContextsFromHelmArchive(t *testing.T) {
	lintCtxs, err := CreateContexts(chartTarball)
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
