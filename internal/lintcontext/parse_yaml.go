package lintcontext

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/internal/k8sutil"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/engine"
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

// LoadObjectsFromPath loads the objects in the file or directory given by `path`
// into the lint context.
func (l *LintContext) LoadObjectsFromPath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return l.loadObjectsFromYAMLFile(path, info)
	}
	err = filepath.Walk(path, func(currentPath string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if isHelm, _ := chartutil.IsChartDir(currentPath); isHelm {
			if err := l.loadObjectsFromHelmChart(currentPath); err != nil {
				return err
			}
			return filepath.SkipDir
		}
		if info.IsDir() {
			return nil
		}
		if !knownYAMLExtensions.Contains(strings.ToLower(filepath.Ext(currentPath))) {
			return nil
		}
		if err := l.loadObjectsFromYAMLFile(currentPath, info); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, "loading from directory %s", path)
	}
	return nil
}

func parseObject(data []byte) (k8sutil.Object, error) {
	obj, _, err := decoder.Decode(data, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode")
	}
	asK8sObj, _ := obj.(k8sutil.Object)
	if asK8sObj == nil {
		return nil, errors.Errorf("object was not a k8s object: %v", obj)
	}
	// TODO: validate
	return asK8sObj, nil
}

type nopWriter struct{}

func (w nopWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (l *LintContext) renderHelmChart(dir string) (map[string]string, error) {
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

func (l *LintContext) loadObjectsFromHelmChart(dir string) error {
	metadata := ObjectMetadata{FilePath: dir}
	renderedFiles, err := l.renderHelmChart(dir)
	if err != nil {
		l.invalidObjects = append(l.invalidObjects, InvalidObject{Metadata: metadata, LoadErr: err})
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

func (l *LintContext) loadObjectFromYAMLReader(filePath string, r *yaml.YAMLReader) error {
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

	obj, err := parseObject(doc)
	if err != nil {
		l.invalidObjects = append(l.invalidObjects, InvalidObject{
			Metadata: metadata,
			LoadErr:  err,
		})
	} else {
		l.objects = append(l.objects, Object{
			Metadata:  metadata,
			K8sObject: obj,
		})
	}
	return nil
}

func (l *LintContext) loadObjectsFromYAMLFile(filePath string, info os.FileInfo) error {
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

func (l *LintContext) loadObjectsFromReader(filePath string, reader io.Reader) error {
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
