package mocks

import (
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
)

// AddContainerToPod adds a mock container to the specified pod under context
func (l *MockLintContext) AddContainerToPod(t *testing.T, podName string, container v1.Container) {
	pod, ok := l.pods[podName]
	require.True(t, ok, "pod with name %s not found", podName)
	// TODO: keep supporting other fields
	pod.Spec.Containers = append(pod.Spec.Containers, container)
}
