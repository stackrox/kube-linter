package mocks

import (
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AddMockPod adds a mock Pod to LintContext
func (l *MockLintContext) AddMockPod(
	podName, namespace, clusterName string,
	labels, annotations map[string]string,
) {
	l.pods[podName] =
		&v1.Pod{
			TypeMeta: metaV1.TypeMeta{},
			ObjectMeta: metaV1.ObjectMeta{
				Name:        podName,
				Namespace:   namespace,
				Labels:      labels,
				Annotations: annotations,
				ClusterName: clusterName,
			},
			Spec:   v1.PodSpec{},
			Status: v1.PodStatus{},
		}
}

// AddSecurityContextToPod adds a security context to the pod specified by name
func (l *MockLintContext) AddSecurityContextToPod(
	podName string,
	runAsUser *int64,
	runAsNonRoot *bool,
) error {
	pod, ok := l.pods[podName]
	if !ok {
		return errors.Errorf("pod with name %q is not found", podName)
	}
	// TODO: keep supporting other fields
	pod.Spec.SecurityContext = &v1.PodSecurityContext{
		SELinuxOptions:      nil,
		WindowsOptions:      nil,
		RunAsUser:           runAsUser,
		RunAsGroup:          nil,
		RunAsNonRoot:        runAsNonRoot,
		SupplementalGroups:  nil,
		FSGroup:             nil,
		Sysctls:             nil,
		FSGroupChangePolicy: nil,
		SeccompProfile:      nil,
	}
	return nil
}
