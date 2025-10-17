package mocks

import (
	"golang.stackrox.io/kube-linter/pkg/k8sutil"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
)

// MockLintContext is mock implementation of the LintContext used in unit tests
type MockLintContext struct {
	objects    map[string]k8sutil.Object
	rawObjects map[string][]byte
}

// Objects returns all the objects under this MockLintContext
func (l *MockLintContext) Objects() []lintcontext.Object {
	result := make([]lintcontext.Object, 0, len(l.objects))
	for key, p := range l.objects {
		metadata := lintcontext.ObjectMetadata{}
		if raw, ok := l.rawObjects[key]; ok {
			metadata.Raw = raw
		}
		result = append(result, lintcontext.Object{Metadata: metadata, K8sObject: p})
	}
	return result
}

// InvalidObjects is not implemented. For now we don't care about invalid objects for mock context.
func (l *MockLintContext) InvalidObjects() []lintcontext.InvalidObject {
	return nil
}

// NewMockContext returns an empty mockLintContext
func NewMockContext() *MockLintContext {
	return &MockLintContext{
		objects:    make(map[string]k8sutil.Object),
		rawObjects: make(map[string][]byte),
	}
}

// AddObject adds an object to the MockLintContext
func (l *MockLintContext) AddObject(key string, obj k8sutil.Object) {
	l.objects[key] = obj
}

// AddObjectWithRaw adds an object to the MockLintContext with raw YAML data
func (l *MockLintContext) AddObjectWithRaw(key string, obj k8sutil.Object, raw []byte) {
	l.objects[key] = obj
	l.rawObjects[key] = raw
}
