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

	"github.com/bmatcuk/doublestar/v4"
	y "github.com/ghodss/yaml"
	kedaV1Alpha1 "github.com/kedacore/keda/v2/apis/keda/v1alpha1"
	ocsAppsV1 "github.com/openshift/api/apps/v1"
	ocpSecV1 "github.com/openshift/api/security/v1"
	"github.com/pkg/errors"
	k8sMonitoring "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"golang.stackrox.io/kube-linter/pkg/k8sutil"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/engine"
	autoscalingV2Beta1 "k8s.io/api/autoscaling/v2beta1"
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

var (
	decoder runtime.Decoder
)

func init() {
	clientScheme := scheme.Scheme

	// Add OpenShift and Autoscaling schema
	schemeBuilder := runtime.NewSchemeBuilder(ocsAppsV1.AddToScheme, autoscalingV2Beta1.AddToScheme, k8sMonitoring.AddToScheme, ocpSecV1.AddToScheme, kedaV1Alpha1.AddToScheme)
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

func (l *lintContextImpl) loadObjectsFromHelmChart(dir string, ignorePaths []string) error {
	renderedFiles, err := l.renderHelmChart(dir)
	if err != nil {
		l.addInvalidObjects(InvalidObject{Metadata: ObjectMetadata{FilePath: dir}, LoadErr: err})
		return nil
	}

	// Paths returned by helm include redundant directory in front, therefore we strip it out.
	return l.loadHelmRenderedTemplates(dir, normalizeDirectoryPaths(renderedFiles), ignorePaths)
}

func (l *lintContextImpl) loadObjectsFromTgzHelmChart(tgzFile string, ignorePaths []string) error {
	renderedFiles, err := l.renderTgzHelmChart(tgzFile)
	if err != nil {
		l.addInvalidObjects(InvalidObject{Metadata: ObjectMetadata{FilePath: tgzFile}, LoadErr: err})
		return nil
	}
	return l.loadHelmRenderedTemplates(tgzFile, renderedFiles, ignorePaths)
}

func (l *lintContextImpl) renderTgzHelmChart(tgzFile string) (map[string]string, error) {
	log.SetOutput(nopWriter{})
	defer log.SetOutput(os.Stderr)

	chrt, err := loader.LoadFile(tgzFile)
	if err != nil {
		return nil, err
	}

	return l.renderChart(tgzFile, chrt)
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
	file, err := os.Open(filepath.Clean(filePath))
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

func (l *lintContextImpl) renderChart(fileName string, chart *chart.Chart) (map[string]string, error) {
	if err := chart.Validate(); err != nil {
		return nil, err
	}

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
		return nil, errors.Wrapf(err, "failed to parse values file %s", indexName)
	}

	return l.renderValues(chart, values)
}

func (l *lintContextImpl) renderTgzHelmChartReader(fileName string, tgzReader io.Reader) (map[string]string, error) {
	// Helm doesn't have great logging behaviour, and can spam stderr, so silence their logging.
	log.SetOutput(nopWriter{})
	defer log.SetOutput(os.Stderr)

	chrt, err := loader.LoadArchive(tgzReader)
	if err != nil {
		return nil, err
	}

	return l.renderChart(fileName, chrt)
}

func (l *lintContextImpl) readObjectsFromTgzHelmChart(fileName string, tgzReader io.Reader, ignoredPaths []string) error {
	renderedFiles, err := l.renderTgzHelmChartReader(fileName, tgzReader)
	if err != nil {
		l.addInvalidObjects(InvalidObject{Metadata: ObjectMetadata{FilePath: fileName}, LoadErr: err})
		return nil
	}
	return l.loadHelmRenderedTemplates(fileName, renderedFiles, ignoredPaths)
}

func (l *lintContextImpl) loadHelmRenderedTemplates(chartPath string, renderedFiles map[string]string, ignorePaths []string) error {
nextFile:
	for path, contents := range renderedFiles {
		pathToTemplate := filepath.Join(chartPath, path)

		for _, path := range ignorePaths {
			ignoreMatch, err := doublestar.PathMatch(path, pathToTemplate)
			if err != nil {
				return errors.Wrapf(err, "could not match pattern %s", path)
			}
			if ignoreMatch {
				continue nextFile
			}
		}

		// Skip NOTES.txt file that may be present among templates but is not a kubernetes resource.
		if strings.HasSuffix(pathToTemplate, string(filepath.Separator)+chartutil.NotesName) {
			continue
		}

		if err := l.loadObjectsFromReader(pathToTemplate, strings.NewReader(contents)); err != nil {
			loadErr := errors.Wrapf(err, "loading object %s from rendered helm chart %s", pathToTemplate, chartPath)
			l.addInvalidObjects(InvalidObject{Metadata: ObjectMetadata{FilePath: pathToTemplate}, LoadErr: loadErr})
		}
	}

	return nil
}

// normalizeDirectoryPaths removes the first element of the path that gets added by the Helm library.
// Helm adds chart name as the first path component, however this is not always correct, e.g. in case the helm chart
// directory was renamed, as shown in https://github.com/stackrox/kube-linter/issues/212
// The function converts mychart/templates/deployment.yaml to templates/deployment.yaml.
func normalizeDirectoryPaths(renderedFiles map[string]string) map[string]string {
	normalizedFiles := make(map[string]string, len(renderedFiles))
	for key, val := range renderedFiles {
		// Go does not seem to have a library function that allows to split the first element of path, therefore
		// splitting "by hand" on path separator char, which is ok if you check path.Split() implementation ;-)
		splitPath := strings.SplitN(key, string(os.PathSeparator), 2)
		if len(splitPath) > 1 {
			normalizedFiles[splitPath[1]] = val
		} else {
			normalizedFiles[splitPath[0]] = val
		}
	}
	return normalizedFiles
}
