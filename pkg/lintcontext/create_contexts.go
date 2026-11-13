package lintcontext

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.stackrox.io/kube-linter/pkg/pathutil"

	"github.com/bmatcuk/doublestar/v4"
	"golang.stackrox.io/kube-linter/internal/set"
	"helm.sh/helm/v3/pkg/chartutil"
	"k8s.io/apimachinery/pkg/runtime"
)

// ReadFromStdin is a path used to indicate reading from os.Stdin
const ReadFromStdin = "-"

var (
	knownYAMLExtensions    = set.NewFrozenStringSet(".yaml", ".yml")
	kustomizationFileNames = []string{"kustomization.yaml", "kustomization.yml"}
)

// Options represent values that can be provided to modify how objects are parsed to create lint contexts
type Options struct {
	// CustomDecoder allows users to supply a non-default decoder to parse k8s objects. This can be used
	// to allow the linter to create contexts for k8s custom resources
	CustomDecoder runtime.Decoder

	// GroupDirectories puts the plain YAML files of all directories into one context instead of
	// one context per directory, so that objects can reference each other across directories
	// (e.g. a Service in services/ selecting a Deployment in deployments/).
	// Helm charts and Kustomize directories are still linted as a context of their own.
	GroupDirectories bool
}

// groupedContextKey is the key of the context shared by all plain YAML files when
// Options.GroupDirectories is set. It can't clash with a path, as paths never contain a NUL byte.
const groupedContextKey = "\x00grouped"

// CreateContexts creates a context. Each context contains a set of files that should be linted
// as a group.
// Each directory of Kube YAML files (or Helm charts) is treated as a separate context.
// Use CreateContextsWithOptions with Options.GroupDirectories to treat the YAML files of all
// directories as being in the same context.
func CreateContexts(ignorePaths []string, filesOrDirs ...string) ([]LintContext, error) {
	return CreateContextsWithOptions(Options{}, ignorePaths, filesOrDirs...)
}

// CreateContextsWithOptions creates a context with additional Options
func CreateContextsWithOptions(options Options, ignorePaths []string, filesOrDirs ...string) ([]LintContext, error) {
	contextsByDir := make(map[string]*lintContextImpl)
fileOrDirsLoop:
	for _, fileOrDir := range filesOrDirs {
		if fileOrDir == ReadFromStdin {
			if _, alreadyExists := contextsByDir[ReadFromStdin]; alreadyExists {
				continue
			}
			ctx := newCtx(options)
			if err := ctx.loadObjectsFromReader("<standard input>", os.Stdin); err != nil {
				return nil, err
			}
			contextsByDir[ReadFromStdin] = ctx
			continue
		}

		for _, path := range ignorePaths {
			// Using doublestar to enable **
			// See https://github.com/golang/go/issues/11862
			globMatch, err := doublestar.PathMatch(path, fileOrDir)
			if err != nil {
				return nil, fmt.Errorf("could not match pattern %s: %w", path, err)
			}
			if globMatch {
				continue fileOrDirsLoop
			}
		}

		err := filepath.Walk(fileOrDir, func(currentPath string, info os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}

			if _, exists := contextsByDir[currentPath]; exists {
				return nil
			}

			for _, path := range ignorePaths {
				absPath, err := pathutil.GetAbsolutPath(currentPath)
				if err != nil {
					return fmt.Errorf("could not get absolute path for %s: %w", currentPath, err)
				}
				globMatch, err := doublestar.PathMatch(path, absPath)
				if err != nil {
					return fmt.Errorf("could not match pattern %s: %w", path, err)
				}
				if globMatch {
					return nil
				}
			}

			if !info.IsDir() {
				if strings.HasSuffix(strings.ToLower(currentPath), ".tgz") {
					ctx := newCtx(options)
					if err := ctx.loadObjectsFromTgzHelmChart(currentPath, ignorePaths); err != nil {
						return fmt.Errorf("loading helm chart %s: %w", currentPath, err)
					}

					contextsByDir[currentPath] = ctx
					return nil
				}

				dirName := filepath.Dir(currentPath)
				if options.GroupDirectories {
					dirName = groupedContextKey
				}
				// Load a file only if it ends in .yaml, OR it was explicitly passed by the user.
				if knownYAMLExtensions.Contains(strings.ToLower(filepath.Ext(currentPath))) || fileOrDir == currentPath {
					ctx := contextsByDir[dirName]
					if ctx == nil {
						ctx = newCtx(options)
						contextsByDir[dirName] = ctx
					}
					if err := ctx.loadObjectsFromYAMLFile(currentPath, info); err != nil {
						return err
					}
				}
				return nil
			}
			if isHelm, _ := chartutil.IsChartDir(currentPath); isHelm {
				// Path has already been loaded, possibly through another argument. Skip.
				if _, alreadyExists := contextsByDir[currentPath]; alreadyExists {
					return nil
				}
				ctx := newCtx(options)
				contextsByDir[currentPath] = ctx
				if err := ctx.loadObjectsFromHelmChart(currentPath, ignorePaths); err != nil {
					return fmt.Errorf("loading helm chart: %w", err)
				}
				return filepath.SkipDir
			}
			if isKustomizeDir(currentPath) {
				// Path has already been loaded, possibly through another argument. Skip.
				if _, alreadyExists := contextsByDir[currentPath]; alreadyExists {
					return nil
				}
				ctx := newCtx(options)
				contextsByDir[currentPath] = ctx
				ctx.loadObjectsFromKustomize(currentPath)
				return filepath.SkipDir
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("loading from path %q: %w", fileOrDir, err)
		}
	}
	dirs := make([]string, 0, len(contextsByDir))
	for dir := range contextsByDir {
		dirs = append(dirs, dir)
	}
	sort.Strings(dirs)
	var contexts []LintContext
	for _, dir := range dirs {
		contexts = append(contexts, contextsByDir[dir])
	}
	return contexts, nil
}

// CreateContextsFromHelmArchive creates a context from TGZ reader of Helm Chart.
// Note: although this function is not used in CLI, it is exposed from kube-linter library and therefore should stay.
// See https://github.com/stackrox/kube-linter/pull/173
func CreateContextsFromHelmArchive(ignorePaths []string, fileName string, tgzReader io.Reader) ([]LintContext, error) {
	ctx := newCtx(Options{})
	if err := ctx.readObjectsFromTgzHelmChart(fileName, tgzReader, ignorePaths); err != nil {
		return nil, err
	}

	return []LintContext{ctx}, nil
}

// isKustomizeDir checks if the given directory contains a kustomization.yaml or kustomization.yml file.
func isKustomizeDir(dirName string) bool {
	for _, fileName := range kustomizationFileNames {
		kustomizationPath := filepath.Join(dirName, fileName)
		if info, err := os.Stat(kustomizationPath); err == nil && !info.IsDir() {
			return true
		}
	}
	return false
}
