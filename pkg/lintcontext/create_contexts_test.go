package lintcontext

import (
	"fmt"
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
	renamedTarball  = "../../tests/testdata/my-renamed-chart-0.1.0.tgz"
	renamedChartDir = "../../tests/testdata/my-renamed-chart"
	mockPath        = "mock path"
)

func TestCreateContextsObjectPaths(t *testing.T) {
	bools := []bool{false, true}

	for _, useTarball := range bools {
		for _, absolute := range bools {
			for _, rename := range bools {
				for _, useFromArchiveFunction := range bools {
					// CreateContextsFromHelmArchive can only be used with tarballs
					if useFromArchiveFunction && !useTarball {
						continue
					}

					testName := fmt.Sprintf("tarball %t, absolute path %t, rename %t, use from archive function %t", useTarball, absolute, rename, useFromArchiveFunction)
					t.Run(testName, func(t *testing.T) {
						createContextsAndVerifyPaths(t, useTarball, absolute, rename, useFromArchiveFunction)
					})
				}
			}
		}
	}
}

func createContextsAndVerifyPaths(t *testing.T, useTarball, useAbsolutePath, rename, useFromArchiveFunction bool) {
	var err error

	// Arrange
	relativePath := map[bool]string{false: chartDirectory, true: chartTarball}[useTarball]
	renamedPath := map[bool]string{false: renamedChartDir, true: renamedTarball}[useTarball]

	testPath := relativePath

	if rename {
		testPath = renamedPath
		require.NoError(t, os.Rename(relativePath, renamedPath))
		defer func() {
			assert.NoError(t, os.Rename(renamedPath, relativePath))
		}()
	}

	if useAbsolutePath {
		testPath, err = filepath.Abs(testPath)
		require.NoError(t, err)
	}

	// Act. The code actually tests either of functions: CreateContextsFromHelmArchive and CreateContexts
	var lintCtxs []LintContext
	if useFromArchiveFunction {
		var file *os.File
		file, err = os.Open(filepath.Clean(testPath))
		require.NoError(t, err)

		defer func() {
			require.NoError(t, file.Close())
		}()

		lintCtxs, err = CreateContextsFromHelmArchive(mockPath, file)
	} else {
		lintCtxs, err = CreateContexts(testPath)
	}
	require.NoError(t, err)

	// Assert
	expectedPath := testPath
	if useTarball {
		expectedPath = path.Join(expectedPath, "mychart")
	}
	if useFromArchiveFunction {
		expectedPath = path.Join(mockPath, "mychart")
	}
	checkObjectPaths(t, verifyAndGetContext(t, lintCtxs).Objects(), expectedPath)
}

func verifyAndGetContext(t *testing.T, lintCtxs []LintContext) LintContext {
	assert.Len(t, lintCtxs, 1, "expecting single lint context to be present")
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
		path.Join(expectedPrefix, "templates/hpa.yaml"),
		path.Join(expectedPrefix, "templates/service.yaml"),
		path.Join(expectedPrefix, "templates/serviceaccount.yaml"),
		path.Join(expectedPrefix, "templates/tests/test-connection.yaml"),
	}
	assert.ElementsMatchf(t, expectedPaths, actualPaths, "expected and actual template paths don't match")
}
