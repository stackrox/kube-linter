package lintcontext

import (
	"encoding/json"

	"golang.stackrox.io/kube-linter/internal/stringutils"
	"golang.stackrox.io/kube-linter/pkg/k8sutil"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// ObjectMetadata is metadata about an object.
type ObjectMetadata struct {
	FilePath string
	Raw      []byte `json:"-"`
}

// An Object references an object that is loaded from a YAML file.
type Object struct {
	Metadata  ObjectMetadata
	K8sObject k8sutil.Object `json:"-"`
}

// K8sObjectInfo contains identifying information about k8s object.
type K8sObjectInfo struct {
	Namespace, Name  string
	GroupVersionKind schema.GroupVersionKind
}

// GetK8sObjectName extracts K8sObjectInfo from Object.K8sObject.
func (o *Object) GetK8sObjectName() K8sObjectInfo {
	return K8sObjectInfo{
		Namespace:        o.K8sObject.GetNamespace(),
		Name:             o.K8sObject.GetName(),
		GroupVersionKind: o.K8sObject.GetObjectKind().GroupVersionKind(),
	}
}

// String provides plain-text representation of k8s object name.
func (n K8sObjectInfo) String() string {
	ns := stringutils.OrDefault(n.Namespace, "<no namespace>")
	return ns + "/" + n.Name + " " + n.GroupVersionKind.String()
}

// MarshalJSON provides custom serialization for Object.
// Object.K8sObject is not serialized directly because that would be too much data. This function limits output to only
// K8sObjectInfo returned for K8sObject.
func (o *Object) MarshalJSON() ([]byte, error) {
	// AliasedObject allows including all Object data without running MarshalJSON (this same function) on it in
	// an infinite loop.
	type AliasedObject Object
	return json.Marshal(&struct {
		*AliasedObject
		K8sObject K8sObjectInfo
	}{
		AliasedObject: (*AliasedObject)(o),
		K8sObject:     o.GetK8sObjectName(),
	})
}

// Check that *Object implements json.Marshaler interface.
var _ json.Marshaler = (*Object)(nil)

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
