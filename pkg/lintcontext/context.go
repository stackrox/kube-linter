package lintcontext

import (
	"golang.stackrox.io/kube-linter/internal/k8sutil"
	"k8s.io/apimachinery/pkg/runtime"
)

// ObjectMetadata is metadata about an object.
type ObjectMetadata struct {
	FilePath string
	Raw      []byte `json:"-"`
}

// An Object references an object that is loaded from a YAML file.
type Object struct {
	Metadata  ObjectMetadata
	K8sObject k8sutil.Object `json:"-"` // TODO: find a way to include sufficient information
}

// An InvalidObject represents something that couldn't be loaded from a YAML file.
type InvalidObject struct {
	Metadata ObjectMetadata
	LoadErr  error
}

// A LintContext represents the context for a lint run.
type LintContext interface {
	Objects() []Object
	InvalidObjects() []InvalidObject
}

type lintContextImpl struct {
	objects        []Object
	invalidObjects []InvalidObject

	customDecoder runtime.Decoder
}

// Objects returns the (valid) objects loaded from this LintContext.
func (l *lintContextImpl) Objects() []Object {
	return l.objects
}

// addObject adds a valid object to this LintContext
func (l *lintContextImpl) addObjects(objs ...Object) {
	l.objects = append(l.objects, objs...)
}

// InvalidObjects returns any objects that we attempted to load, but which were invalid.
func (l *lintContextImpl) InvalidObjects() []InvalidObject {
	return l.invalidObjects
}

// addInvalidObject adds an invalid object to this LintContext
func (l *lintContextImpl) addInvalidObjects(objs ...InvalidObject) {
	l.invalidObjects = append(l.invalidObjects, objs...)
}

// new returns a ready-to-use, empty, lintContextImpl.
func newCtx(options Options) *lintContextImpl {
	return &lintContextImpl{
		customDecoder: options.CustomDecoder,
	}
}
