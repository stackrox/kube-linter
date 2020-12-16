package lintcontext

import (
	"golang.stackrox.io/kube-linter/internal/k8sutil"
)

// ObjectMetadata is metadata about an object.
type ObjectMetadata struct {
	FilePath string
	Raw      []byte
}

// An Object references an object that is loaded from a YAML file.
type Object struct {
	Metadata  ObjectMetadata
	K8sObject k8sutil.Object
}

// An InvalidObject represents something that couldn't be loaded from a YAML file.
type InvalidObject struct {
	Metadata ObjectMetadata
	LoadErr  error
}

// A LintContext represents the context for a lint run.
type LintContext interface {
	GetObjects() []Object
	GetInvalidObjects() []InvalidObject
}

type lintContextImpl struct {
	objects        []Object
	invalidObjects []InvalidObject
}

// Objects returns the (valid) objects loaded from this LintContext.
func (l *lintContextImpl) GetObjects() []Object {
	return l.objects
}

// AddObject adds a valid object to this LintContext
func (l *lintContextImpl) addObjects(objs ...Object) {
	l.objects = append(l.objects, objs...)
}

// InvalidObjects returns any objects that we attempted to load, but which were invalid.
func (l *lintContextImpl) GetInvalidObjects() []InvalidObject {
	return l.invalidObjects
}

// AddInvalidObject adds an invalid object to this LintContext
func (l *lintContextImpl) addInvalidObjects(objs ...InvalidObject) {
	l.invalidObjects = append(l.invalidObjects, objs...)
}

// New returns a ready-to-use, empty, lint context.
func New() *lintContextImpl {
	return &lintContextImpl{}
}
