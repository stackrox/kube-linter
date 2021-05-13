package mocks

import (
	"golang.stackrox.io/kube-linter/pkg/k8sutil"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
)

// MockLintContext is mock implementation of the LintContext used in unit tests
type MockLintContext struct {
	objects map[string]k8sutil.Object
}

// Objects returns all the objects under this MockLintContext
func (l *MockLintContext) Objects() []lintcontext.Object {
	result := make([]lintcontext.Object, 0, len(l.objects))
	for _, p := range l.objects {
		result = append(result, lintcontext.Object{Metadata: lintcontext.ObjectMetadata{}, K8sObject: p})
	}
	return result
}

// InvalidObjects is not implemented. For now we don't care about invalid objects for mock context.
func (l *MockLintContext) InvalidObjects() []lintcontext.InvalidObject {
	return nil
}

// NewMockContext returns an empty mockLintContext
func NewMockContext() *MockLintContext {
	return &MockLintContext{objects: make(map[string]k8sutil.Object)}
}
