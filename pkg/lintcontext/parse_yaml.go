package lintcontext

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	y "github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/internal/k8sutil"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/engine"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
)

const (
	// The max file size, in bytes, that we will load.
	// TODO: make it configurable.
	maxFileSizeBytes = 10 * 1024 * 1024
)

var (
	clientSchema = scheme.Scheme
	decoder      = serializer.NewCodecFactory(clientSchema).UniversalDeserializer()
)

func parseObjects(data []byte) ([]k8sutil.Object, error) {
	obj, _, err := decoder.Decode(data, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode")
	}
	if list, ok := obj.(*v1.List); ok {
		objs := make([]k8sutil.Object, 0, len(list.Items))
		for i, item := range list.Items {
			obj, _, err := decoder.Decode(item.Raw, nil, nil)
			if err != nil {
				return nil, errors.Wrapf(err, "decoding item %d in the list", i)
			}
			asK8sObj, _ := obj.(k8sutil.Object)
			if asK8sObj == nil {
				return nil, errors.Errorf("object was not a k8s object: %v", obj)
			}
			objs = append(objs, asK8sObj)
		}
		return objs, nil
	}
	asK8sObj, _ := obj.(k8sutil.Object)
	if asK8sObj == nil {
		return nil, errors.Errorf("object was not a k8s object: %v", obj)
	}
	// TODO: validate
	return []k8sutil.Object{asK8sObj}, nil
}

type nopWriter struct{}

func (w nopWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (l *lintContextImpl) renderHelmChart(dir string) (map[string]string, error) {
	// Helm doesn't have great logging behaviour, and can spam stderr, so silence their logging.
	// TODO: capture these logs.
	log.SetOutput(nopWriter{})
	defer log.SetOutput(os.Stderr)
	chrt, err := loader.Load(dir)
	if err != nil {
		return nil, err
	}
	if err := chrt.Validate(); err != nil {
		return nil, err
	}
	valOpts := &values.Options{ValueFiles: []string{filepath.Join(dir, "values.yaml")}}
	values, err := valOpts.MergeValues(nil)
	if err != nil {
		return nil, errors.Wrap(err, "loading values.yaml file")
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

func (l *lintContextImpl) loadObjectsFromHelmChart(dir string) error {
	metadata := ObjectMetadata{FilePath: dir}
	renderedFiles, err := l.renderHelmChart(dir)
	if err != nil {
		l.addInvalidObjects(InvalidObject{Metadata: metadata, LoadErr: err})
		return nil
	}
	for path, contents := range renderedFiles {
		// The first element of path will be the same as the last element of dir, because
		// Helm duplicates it.
		pathToTemplate := filepath.Join(filepath.Dir(dir), path)
		if err := l.loadObjectsFromReader(pathToTemplate, strings.NewReader(contents)); err != nil {
			return errors.Wrapf(err, "loading objects from rendered helm chart %s/%s", dir, pathToTemplate)
		}
	}
	return nil
}

func (l *lintContextImpl) loadObjectsFromTgzHelmChart(tgzFile string) error {
	metadata := ObjectMetadata{FilePath: tgzFile}
	renderedFiles, err := l.renderTgzHelmChart(tgzFile)
	if err != nil {
		l.invalidObjects = append(l.invalidObjects, InvalidObject{Metadata: metadata, LoadErr: err})
		return nil
	}
	for path, contents := range renderedFiles {
		// The first element of path will be the same as the last element of tgzFile, because
		// Helm duplicates it.
		pathToTemplate := filepath.Join(filepath.Dir(tgzFile), path)
		if err := l.loadObjectsFromReader(pathToTemplate, strings.NewReader(contents)); err != nil {
			return errors.Wrapf(err, "loading objects from rendered helm chart %s/%s", tgzFile, pathToTemplate)
		}
	}
	return nil
}

func (l *lintContextImpl) renderTgzHelmChart(tgzFile string) (map[string]string, error) {
	log.SetOutput(nopWriter{})
	defer log.SetOutput(os.Stderr)
	chrt, err := loader.LoadFile(tgzFile)

	if err != nil {
		return nil, err
	}
	if err := chrt.Validate(); err != nil {
		return nil, err
	}

	valuesIndex := -1
	for i, f := range chrt.Raw {
		if f.Name == "values.yaml" {
			valuesIndex = i
			break
		}
	}

	indexName := filepath.Join(tgzFile, "values.yaml")
	if valuesIndex == -1 {
		return nil, errors.Errorf("%s not found", indexName)
	}

	values, err := l.parseValues(indexName, chrt.Raw[valuesIndex].Data)
	if err != nil {
		return nil, errors.Wrap(err, "loading values.yaml file")
	}

	return l.renderValues(chrt, values)
}

func (l *lintContextImpl) parseValues(filePath string, bytes []byte) (map[string]interface{}, error) {
	currentMap := map[string]interface{}{}

	if err := y.Unmarshal(bytes, &currentMap); err != nil {
		return nil, errors.Wrapf(err, "failed to parse %s", filePath)
	}

	return currentMap, nil
}

func (l *lintContextImpl) loadObjectFromYAMLReader(filePath string, r *yaml.YAMLReader) error {
	doc, err := r.Read()
	if err != nil {
		return err
	}
	doc = bytes.TrimSpace(doc)
	if len(doc) == 0 {
		return nil
	}

	metadata := ObjectMetadata{
		FilePath: filePath,
		Raw:      doc,
	}

	objs, err := parseObjects(doc)
	if err != nil {
		l.addInvalidObjects(InvalidObject{
			Metadata: metadata,
			LoadErr:  err,
		})
		return nil
	}
	for _, obj := range objs {
		l.addObjects(Object{
			Metadata:  metadata,
			K8sObject: obj,
		})
	}
	return nil
}

func (l *lintContextImpl) loadObjectsFromYAMLFile(filePath string, info os.FileInfo) error {
	if info.Size() > maxFileSizeBytes {
		return nil
	}
	file, err := os.Open(filePath)
	if err != nil {
		return errors.Wrapf(err, "opening file at %s", filePath)
	}
	defer func() {
		_ = file.Close()
	}()

	return l.loadObjectsFromReader(filePath, file)
}

func (l *lintContextImpl) loadObjectsFromReader(filePath string, reader io.Reader) error {
	yamlReader := yaml.NewYAMLReader(bufio.NewReader(reader))
	for {
		if err := l.loadObjectFromYAMLReader(filePath, yamlReader); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}

}
