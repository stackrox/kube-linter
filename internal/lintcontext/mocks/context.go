package mocks

import (
	"golang.stackrox.io/kube-linter/internal/lintcontext"
	v1 "k8s.io/api/core/v1"
)

// MockLintContext is mock implementation of the LintContext used in unit tests
type MockLintContext struct {
	pods map[string]*v1.Pod
}

// GetObjects returns all the objects under this MockLintContext
func (l *MockLintContext) GetObjects() []lintcontext.Object {
	result := make([]lintcontext.Object, 0, len(l.pods))
	for _, p := range l.pods {
		result = append(result, lintcontext.Object{Metadata: lintcontext.ObjectMetadata{}, K8sObject: p})
	}
	return result
}

// GetInvalidObjects is not implemented. For now we don't care about invalid objects for mock context.
func (l *MockLintContext) GetInvalidObjects() []lintcontext.InvalidObject {
	return nil
}

// NewMockContext returns an empty mockLintContext
func NewMockContext() *MockLintContext {
	return &MockLintContext{pods: make(map[string]*v1.Pod)}
}
