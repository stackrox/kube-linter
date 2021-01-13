package mocks

import (
	"testing"

	"github.com/stretchr/testify/require"
	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

// AddContainerToDeployment adds a mock container to the specified pod under context
func (l *MockLintContext) AddContainerToDeployment(t *testing.T, deploymentName string, container v1.Container) {
	deployment, ok := l.objects[deploymentName].(*appsV1.Deployment)
	require.True(t, ok, "deployment with name %s not found", deploymentName)
	// TODO: keep supporting other fields
	deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, container)
}
