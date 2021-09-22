package lintcontext

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	chartTarball    = "../../tests/testdata/mychart-0.1.0.tgz"
	chartDirectory  = "../../tests/testdata/mychart"
	renamedChartDir = "../../tests/testdata/my-renamed-chart"
)

func TestCreateContextFromHelmArchive(t *testing.T) {
	lintCtx := createSingleContext(t, chartTarball)

	checkObjectPaths(t, lintCtx.Objects(), path.Join(chartTarball, "mychart"))
}

func TestCreateContextFromHelmArchiveAbsolutePath(t *testing.T) {
	absPath, err := filepath.Abs(chartTarball)
	require.NoError(t, err)

	lintCtx := createSingleContext(t, absPath)

	checkObjectPaths(t, lintCtx.Objects(), path.Join(absPath, "mychart"))
}

func TestCreateContextFromHelmDirectory(t *testing.T) {
	lintCtx := createSingleContext(t, chartDirectory)

	checkObjectPaths(t, lintCtx.Objects(), chartDirectory)
}

func TestCreateContextFromHelmDirectoryAbsolutePath(t *testing.T) {
	absPath, err := filepath.Abs(chartDirectory)
	require.NoError(t, err)
	lintCtx := createSingleContext(t, absPath)

	checkObjectPaths(t, lintCtx.Objects(), absPath)
}

func TestCreateContextFromRenamedHelmDirectory(t *testing.T) {
	require.NoError(t, os.Rename(chartDirectory, renamedChartDir))
	defer func() {
		assert.NoError(t, os.Rename(renamedChartDir, chartDirectory))
	}()

	lintCtx := createSingleContext(t, renamedChartDir)

	checkObjectPaths(t, lintCtx.Objects(), renamedChartDir)
}

func TestCreateContextFromRenamedHelmDirectoryAbsolutePath(t *testing.T) {
	require.NoError(t, os.Rename(chartDirectory, renamedChartDir))
	defer func() {
		assert.NoError(t, os.Rename(renamedChartDir, chartDirectory))
	}()

	absPath, err := filepath.Abs(renamedChartDir)
	require.NoError(t, err)

	lintCtx := createSingleContext(t, absPath)

	checkObjectPaths(t, lintCtx.Objects(), absPath)
}

func createSingleContext(t *testing.T, path string) LintContext {
	lintCtxs, err := CreateContexts(path)
	require.NoError(t, err)
	assert.Len(t, lintCtxs, 1, "expecting single lint context to be returned")

	lintCtx := lintCtxs[0]

	assert.NotEmpty(t, lintCtx.Objects(), "no valid objects found")
	assert.Empty(t, lintCtx.InvalidObjects(), "no invalid objects expected")

	return lintCtx
}

func checkObjectPaths(t *testing.T, objects []Object, expectedPrefix string) {
	actualPaths := make([]string, 0, len(objects))
	for _, obj := range objects {
		actualPaths = append(actualPaths, obj.Metadata.FilePath)
	}
	expectedPaths := []string{
		path.Join(expectedPrefix, "templates/deployment.yaml"),
		path.Join(expectedPrefix, "templates/service.yaml"),
		path.Join(expectedPrefix, "templates/serviceaccount.yaml"),
		path.Join(expectedPrefix, "templates/tests/test-connection.yaml"),
	}
	assert.ElementsMatchf(t, expectedPaths, actualPaths, "expected and actual template paths don't match")
}
