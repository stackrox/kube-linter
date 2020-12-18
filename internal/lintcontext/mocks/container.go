package mocks

import (
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
)

// AddContainerToPod adds a mock container to the specified pod under context
func (l *MockLintContext) AddContainerToPod(
	podName, containerName, image string,
	ports []v1.ContainerPort,
	env []v1.EnvVar,
	sc *v1.SecurityContext,
) error {
	pod, ok := l.pods[podName]
	if !ok {
		return errors.Errorf("pod with name %q is not found", podName)
	}
	// TODO: keep supporting other fields
	pod.Spec.Containers = append(pod.Spec.Containers, v1.Container{
		Name:            containerName,
		Image:           image,
		Ports:           ports,
		Env:             env,
		Resources:       v1.ResourceRequirements{},
		SecurityContext: sc,
	})
	return nil
}
