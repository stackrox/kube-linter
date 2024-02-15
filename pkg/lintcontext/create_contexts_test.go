package lintcontext

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"golang.stackrox.io/kube-linter/pkg/pathutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	chartTarball       = "../../tests/testdata/mychart-0.1.0.tgz"
	chartDirectory     = "../../tests/testdata/mychart"
	renamedTarball     = "../../tests/testdata/my-renamed-chart-0.1.0.tgz"
	renamedChartDir    = "../../tests/testdata/my-renamed-chart"
	mockIgnorePath     = "../../tests/testdata/**"
	mockGlobIgnorePath = "../../tests/**"
	mockPath           = "mock path"
)

func TestCreateContextsWithIgnorePaths(t *testing.T) {
	ignoredPaths := []string{
		"../../.github/**",
		"../../.golangci?yml",
		"../../.goreleaser.yaml",
		"../../.pre-commit-hooks*",
		"../../dist/**/*",
		"../../pkg/**/*",
		"/**/*/checks/**/*",
		"/**/*/test_helper/**/*",
		"/**/*/testdata/**/*",
	}
	ignoredAbsPaths := make([]string, 0, len(ignoredPaths))
	for _, p := range ignoredPaths {
		abs, err := pathutil.GetAbsolutPath(p)
		assert.NoError(t, err)
		ignoredAbsPaths = append(ignoredAbsPaths, abs)
	}

	testPath := "../../"
	contexts, err := CreateContexts(ignoredAbsPaths, testPath)
	assert.NoError(t, err)
	checkEmptyLintContext(t, contexts)
}

func TestIgnoreSubchartManifests(t *testing.T) {
	ignorePaths := []string{
		"../../tests/testdata/mychart/charts/**",
	}
	dir := "../../tests/testdata/mychart"

	lintCtxs, err := CreateContexts(ignorePaths, dir)
	require.NoError(t, err)
	lintCtx := lintCtxs[0]
	objects := lintCtx.Objects()

	actualPaths := make([]string, 0, len(objects))
	for _, obj := range objects {
		actualPaths = append(actualPaths, obj.Metadata.FilePath)
	}

	expectedPaths := []string{
		"../../tests/testdata/mychart/templates/serviceaccount.yaml",
		"../../tests/testdata/mychart/templates/service.yaml",
		"../../tests/testdata/mychart/templates/hpa.yaml",
		"../../tests/testdata/mychart/templates/deployment.yaml",
		"../../tests/testdata/mychart/templates/tests/test-connection.yaml",
	}

	assert.ElementsMatch(t, expectedPaths, actualPaths)
}

func TestCreateContextsObjectPaths(t *testing.T) {
	bools := []bool{false, true}

	for _, useTarball := range bools {
		for _, absolute := range bools {
			for _, rename := range bools {
				for _, useFromArchiveFunction := range bools {
					for _, useGlob := range bools {
						for _, useIgnorePaths := range bools {
							// CreateContextsFromHelmArchive can only be used with tarballs
							if useFromArchiveFunction && !useTarball {
								continue
							}

							testName := fmt.Sprintf("tarball %t, absolute path %t, rename %t, use from archive function %t, ignore paths: %t (use glob: %t)", useTarball, absolute, rename, useFromArchiveFunction, useIgnorePaths, useGlob)
							t.Run(testName, func(t *testing.T) {
								createContextsAndVerifyPaths(t, useTarball, absolute, rename, useFromArchiveFunction, useIgnorePaths, useGlob)
							})
						}
					}
				}
			}
		}
	}
}

func createContextsAndVerifyPaths(t *testing.T, useTarball, useAbsolutePath, rename, useFromArchiveFunction, useIgnorePaths, useGlob bool) {
	var err error

	// Arrange
	relativePath := map[bool]string{false: chartDirectory, true: chartTarball}[useTarball]
	renamedPath := map[bool]string{false: renamedChartDir, true: renamedTarball}[useTarball]

	testPath := relativePath
	testIgnorePath := mockIgnorePath
	testIgnorePaths := make([]string, 0)

	if rename {
		testPath = renamedPath
		require.NoError(t, os.Rename(relativePath, renamedPath))
		defer func() {
			assert.NoError(t, os.Rename(renamedPath, relativePath))
		}()
	}

	if useGlob {
		testIgnorePath = mockGlobIgnorePath
	}

	if useAbsolutePath {
		testPath, err = filepath.Abs(testPath)
		require.NoError(t, err)
		testIgnorePath, err = filepath.Abs(testIgnorePath)
		require.NoError(t, err)
	}

	if useIgnorePaths {
		testIgnorePaths = append(testIgnorePaths, testIgnorePath)
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

		lintCtxs, err = CreateContextsFromHelmArchive(testIgnorePaths, mockPath, file)
	} else {
		lintCtxs, err = CreateContexts(testIgnorePaths, testPath)
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

	// IgnorePaths is only used for non helm cases
	if useIgnorePaths && !useFromArchiveFunction {
		checkEmptyLintContext(t, lintCtxs)
	} else {
		checkObjectPaths(t, verifyAndGetContext(t, lintCtxs).Objects(), expectedPath)
	}
}

func checkEmptyLintContext(t *testing.T, lintCtxs []LintContext) {
	assert.Empty(t, lintCtxs, "expecting no lint context")
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
		path.Join(expectedPrefix, "charts/subchart/templates/deployment.yaml"),
	}
	assert.ElementsMatchf(t, expectedPaths, actualPaths, "expected and actual template paths don't match")
}
