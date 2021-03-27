package lintcontext

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	ocsAppsV1 "github.com/openshift/api/apps/v1"
	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/pkg/k8sutil"
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
