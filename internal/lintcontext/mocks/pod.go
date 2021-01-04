package mocks

import (
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AddMockPod adds a mock Pod to LintContext
func (l *MockLintContext) AddMockPod(
	podName, namespace string,
	labels, annotations map[string]string,
) {
	l.pods[podName] =
		&v1.Pod{
			ObjectMeta: metaV1.ObjectMeta{
				Name:        podName,
				Namespace:   namespace,
				Labels:      labels,
				Annotations: annotations,
			},
		}
}

// AddSecurityContextToPod adds a security context to the pod specified by name
func (l *MockLintContext) AddSecurityContextToPod(t *testing.T, podName string,
	securityContext *v1.PodSecurityContext) {
	pod, ok := l.pods[podName]
	require.True(t, ok, "pod with name %s not added", podName)
	// TODO: keep supporting other fields
	pod.Spec.SecurityContext = securityContext
}
