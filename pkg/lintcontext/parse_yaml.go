package lintcontext

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	y "github.com/ghodss/yaml"
	ocsAppsV1 "github.com/openshift/api/apps/v1"
	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/pkg/k8sutil"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/engine"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
)

const (
	// The max file size, in bytes, that we will load.
	// TODO: make it configurable.
	maxFileSizeBytes = 10 * 1024 * 1024
)

type chartType int

const (
	// Helm chart is set of files in a directory
	chartInDirectory chartType = iota
	// Helm chart is compressed in .tgz archive
	chartInTgzFile
)

var (
	decoder runtime.Decoder
)

func init() {
	clientScheme := scheme.Scheme

	// Add OpenShift schema
	schemeBuilder := runtime.NewSchemeBuilder(ocsAppsV1.AddToScheme)
	if err := schemeBuilder.AddToScheme(clientScheme); err != nil {
		panic(fmt.Sprintf("Can not add OpenShift schema %v", err))
	}
	decoder = serializer.NewCodecFactory(clientScheme).UniversalDeserializer()
}

func parseObjects(data []byte, d runtime.Decoder) ([]k8sutil.Object, error) {
	if d == nil {
		d = decoder
	}
	obj, _, err := d.Decode(data, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode")
	}
	if list, ok := obj.(*v1.List); ok {
		objs := make([]k8sutil.Object, 0, len(list.Items))
		for i, item := range list.Items {
			obj, _, err := d.Decode(item.Raw, nil, nil)
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

	values, err := l.loadTgzHelmValues(tgzFile, chrt)
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

func (l *lintContextImpl) loadObjectsFromHelmChart(chartPath string, chartType chartType) error {
	metadata := ObjectMetadata{FilePath: chartPath}

	var renderedFiles map[string]string
	var err error
	switch chartType {
	case chartInDirectory:
		renderedFiles, err = l.renderHelmChart(chartPath)
	case chartInTgzFile:
		renderedFiles, err = l.renderTgzHelmChart(chartPath)
	default:
		return errors.Errorf("unknown chart type %q", chartType)
	}
	if err != nil {
		l.addInvalidObjects(InvalidObject{Metadata: metadata, LoadErr: err})
		return nil
	}

	for path, contents := range renderedFiles {
		// The first element of path will be the same as the last element of chartPath, because
		// Helm duplicates it.
		pathToTemplate := filepath.Join(filepath.Dir(chartPath), path)
		if err := l.loadObjectsFromReader(pathToTemplate, strings.NewReader(contents)); err != nil {
			return errors.Wrapf(err, "loading objects from rendered helm chart %s/%s", chartPath, pathToTemplate)
		}
	}
	return nil
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

	objs, err := parseObjects(doc, l.customDecoder)
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

func (l *lintContextImpl) loadTgzHelmValues(fileName string, chart *chart.Chart) (map[string]interface{}, error) {
	valuesIndex := -1
	for i, f := range chart.Raw {
		if f.Name == "values.yaml" {
			valuesIndex = i
			break
		}
	}

	indexName := filepath.Join(fileName, "values.yaml")
	if valuesIndex == -1 {
		return nil, errors.Errorf("%s not found", indexName)
	}

	values := map[string]interface{}{}
	if err := y.Unmarshal(chart.Raw[valuesIndex].Data, &values); err != nil {
		return nil, errors.Wrapf(err, "failed to parse %s", indexName)
	}

	return values, nil
}
