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
	chartTarball           = "../../tests/testdata/mychart-0.1.0.tgz"
	chartDirectory         = "../../tests/testdata/mychart"
	renamedTarball         = "../../tests/testdata/my-renamed-chart-0.1.0.tgz"
	renamedChartDir        = "../../tests/testdata/my-renamed-chart"
	kustomizeDirectory     = "../../tests/testdata/mykustomize"
	kustomizeDeprecatedDir = "../../tests/testdata/mykustomize-deprecated"
	mockIgnorePath         = "../../tests/testdata/**"
	mockGlobIgnorePath     = "../../tests/**"
	mockPath               = "mock path"
)

func TestCreateContextsWithIgnorePaths(t *testing.T) {
	ignoredPaths := []string{
		"../../.github/**",
		"../../.golangci?yml",
		"../../.goreleaser.yaml",
		"../../.pre-commit-hooks*",
		"../../dist/**/*",
		"../../pkg/**/*",
		"../../demo/**",
		"../../stackrox-kube-linter-bug-example/**",
		"../../tests/**/*",
		"../../cmd/**/*",
		"../../docs/**/*",
		"../../internal/**/*",
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

func TestKustomizeContextCreation(t *testing.T) {
	// Test that kustomize directory is loaded correctly
	lintCtxs, err := CreateContexts(nil, kustomizeDirectory)
	require.NoError(t, err)
	require.Len(t, lintCtxs, 1, "expecting single lint context to be present")

	lintCtx := lintCtxs[0]
	objects := lintCtx.Objects()

	assert.NotEmpty(t, objects, "expecting kustomize objects to be loaded")
	assert.Empty(t, lintCtx.InvalidObjects(), "no invalid objects expected")

	// Verify that we have the expected objects (Deployment and Service)
	assert.Len(t, objects, 2, "expecting 2 objects from kustomization")

	// Check that objects have the kustomize transformations applied
	foundDeployment := false
	foundService := false
	for _, obj := range objects {
		k8sObj := obj.K8sObject
		if k8sObj.GetObjectKind().GroupVersionKind().Kind == "Deployment" {
			foundDeployment = true
			// Check that kustomize namePrefix was applied
			assert.Contains(t, k8sObj.GetName(), "kustomize-", "deployment should have kustomize namePrefix")
		}
		if k8sObj.GetObjectKind().GroupVersionKind().Kind == "Service" {
			foundService = true
			// Check that kustomize namePrefix was applied
			assert.Contains(t, k8sObj.GetName(), "kustomize-", "service should have kustomize namePrefix")
		}
	}

	assert.True(t, foundDeployment, "deployment should be present")
	assert.True(t, foundService, "service should be present")
}

func TestKustomizeWithIgnorePaths(t *testing.T) {
	// Test that ignore paths work with kustomize
	ignorePaths := []string{kustomizeDirectory + "/**"}
	lintCtxs, err := CreateContexts(ignorePaths, kustomizeDirectory)
	require.NoError(t, err)

	// Should be empty because we ignored the kustomize directory
	checkEmptyLintContext(t, lintCtxs)
}

func TestKustomizeWithDeprecatedSyntax(t *testing.T) {
	// Test that kustomize deprecation warnings are ignored (not converted to errors)
	// Similar to how helm warnings are suppressed
	lintCtxs, err := CreateContexts(nil, kustomizeDeprecatedDir)
	require.NoError(t, err, "CreateContexts should not return an error")
	require.Len(t, lintCtxs, 1, "expecting single lint context to be present")

	lintCtx := lintCtxs[0]

	// Warnings should be ignored, so no invalid objects
	invalidObjects := lintCtx.InvalidObjects()
	assert.Empty(t, invalidObjects, "deprecated syntax should not create invalid objects")

	// The kustomization should still load successfully despite using deprecated syntax
	objects := lintCtx.Objects()
	assert.NotEmpty(t, objects, "objects should be loaded even with deprecated 'bases' field")

	// Verify we have the expected objects from the kustomization
	assert.Len(t, objects, 2, "expecting 2 objects (Deployment and Service) from deprecated kustomization")
}

func TestWarningSuppressionParity(t *testing.T) {
	// This test documents that both Helm and Kustomize suppress warnings
	// Helm: Uses nopWriter{} to discard log output (parse_yaml.go:102-106)
	// Kustomize: Uses WithWarningHandler that returns nil (parse_yaml.go:336-339)

	t.Run("Helm suppresses warnings via nopWriter", func(t *testing.T) {
		// Test that helm charts load without error even if they generate warnings
		// The existing helm tests verify this behavior by successfully loading charts
		lintCtxs, err := CreateContexts(nil, chartDirectory)
		require.NoError(t, err)
		require.Len(t, lintCtxs, 1)

		lintCtx := lintCtxs[0]
		assert.Empty(t, lintCtx.InvalidObjects(), "helm charts should load despite any warnings")
		assert.NotEmpty(t, lintCtx.Objects(), "helm charts should produce objects")
	})

	t.Run("Kustomize suppresses warnings via handler", func(t *testing.T) {
		// Test that kustomizations with deprecation warnings still load
		lintCtxs, err := CreateContexts(nil, kustomizeDeprecatedDir)
		require.NoError(t, err)
		require.Len(t, lintCtxs, 1)

		lintCtx := lintCtxs[0]
		assert.Empty(t, lintCtx.InvalidObjects(), "warnings should be suppressed, not converted to errors")
		assert.NotEmpty(t, lintCtx.Objects(), "objects should load despite warnings")
	})
}
