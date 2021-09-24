package lintcontext

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
)

func (l *lintContextImpl) loadObjectsFromHelmChart(path string, options HelmOptions) error {
	metadata := ObjectMetadata{FilePath: path}
	renderedFiles, err := l.renderHelmChart(path, options)
	if err != nil {
		l.addInvalidObjects(InvalidObject{Metadata: metadata, LoadErr: err})
		return nil
	}
	for path, contents := range renderedFiles {
		// The first element of path will be the same as the last element of dir, because
		// Helm duplicates it.
		pathToTemplate := filepath.Join(filepath.Dir(path), path)
		if err := l.loadObjectsFromReader(pathToTemplate, strings.NewReader(contents)); err != nil {
			return errors.Wrapf(err, "loading objects from rendered helm chart %s/%s", path, pathToTemplate)
		}
	}
	return nil
}

func (l *lintContextImpl) renderHelmChart(path string, options HelmOptions) (map[string]string, error) {
	// Helm doesn't have great logging behaviour, and can spam stderr, so silence their logging.
	// TODO: capture these logs.
	log.SetOutput(nopWriter{})
	defer log.SetOutput(os.Stderr)

	var chrt *chart.Chart
	var err error
	if options.FromDir && options.FromArchive {
		return nil, errors.New("cannot specify that helm chart is both a directory and an archive")
	}

	switch {
	case options.FromArchive:
		chrt, err = loader.LoadFile(path)
	case options.FromDir:
		chrt, err = loader.Load(path)
	default:
		chrt, err = loader.LoadArchive(options.FromReader)
	}
	if err != nil {
		return nil, err
	}

	if err := chrt.Validate(); err != nil {
		return nil, err
	}
	values, err := l.helmValuesOptions.MergeValues(nil)
	if err != nil {
		return nil, errors.Wrap(err, "loading provided Helm value options")
	}

	return l.renderValues(chrt, values)
}

func (l *lintContextImpl) renderValues(chrt *chart.Chart, values map[string]interface{}) (map[string]string, error) {
	valuesToRender, err := chartutil.ToRenderValues(chrt, values, chartutil.ReleaseOptions{Name: "test-release", Namespace: "default"}, nil)
	if err != nil {
		return nil, err
	}

	e := engine.Engine{LintMode: true}
	rendered, err := e.Render(chrt, valuesToRender)
	if err != nil {
		return nil, errors.Wrap(err, "failed to render")
	}

	return rendered, nil
}
