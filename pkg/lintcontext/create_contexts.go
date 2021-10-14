package lintcontext

import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/internal/set"
	"helm.sh/helm/v3/pkg/chartutil"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	knownYAMLExtensions = set.NewFrozenStringSet(".yaml", ".yml")
)

// Options represent values that can be provided to modify how objects are parsed to create lint contexts
type Options struct {
	// CustomDecoder allows users to supply a non-default decoder to parse k8s objects. This can be used
	// to allow the linter to create contexts for k8s custom resources
	CustomDecoder runtime.Decoder
}

// CreateContexts creates a context. Each context contains a set of files that should be linted
// as a group.
// Currently, each directory of Kube YAML files (or Helm charts) are treated as a separate context.
// TODO: Figure out if it's useful to allow people to specify that files spanning different directories
// should be treated as being in the same context.
func CreateContexts(filesOrDirs ...string) ([]LintContext, error) {
	return CreateContextsWithOptions(Options{}, filesOrDirs...)
}

// CreateContextsWithOptions creates a context with additional Options
func CreateContextsWithOptions(options Options, filesOrDirs ...string) ([]LintContext, error) {
	contextsByDir := make(map[string]*lintContextImpl)
	for _, fileOrDir := range filesOrDirs {
		// Stdin
		if fileOrDir == "-" {
			if _, alreadyExists := contextsByDir["-"]; alreadyExists {
				continue
			}
			ctx := newCtx(options)
			if err := ctx.loadObjectsFromReader("<standard input>", os.Stdin); err != nil {
				return nil, err
			}
			contextsByDir["-"] = ctx
			continue
		}

		err := filepath.Walk(fileOrDir, func(currentPath string, info os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}

			if _, exists := contextsByDir[currentPath]; exists {
				return nil
			}

			if !info.IsDir() {
				if strings.HasSuffix(strings.ToLower(currentPath), ".tgz") {
					ctx := newCtx(options)
					ctx.loadObjectsFromTgzHelmChart(currentPath)
					contextsByDir[currentPath] = ctx
					return nil
				}

				dirName := filepath.Dir(currentPath)
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
				ctx.loadObjectsFromHelmChart(currentPath)
				return filepath.SkipDir
			}
			return nil
		})
		if err != nil {
			return nil, errors.Wrapf(err, "loading from path %q", fileOrDir)
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
func CreateContextsFromHelmArchive(fileName string, tgzReader io.Reader) ([]LintContext, error) {
	ctx := newCtx(Options{})
	ctx.readObjectsFromTgzHelmChart(fileName, tgzReader)

	return []LintContext{ctx}, nil
}
