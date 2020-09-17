package lintcontext

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/internal/k8sutil"
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
		if currentPath == path {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(currentPath) != ".yaml" {
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

	yamlReader := yaml.NewYAMLReader(bufio.NewReader(file))
	for {
		if err := l.loadObjectFromYAMLReader(filePath, yamlReader); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

// Objects returns the (valid) objects loaded from this LintContext.
func (l *LintContext) Objects() []Object {
	return l.objects
}

// InvalidObjects returns any objects that we attempted to load, but which were invalid.
func (l *LintContext) InvalidObjects() []InvalidObject {
	return l.invalidObjects
}
