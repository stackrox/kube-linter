package mocks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

func TestAddContainersToDeployment(t *testing.T) {
	ctx := NewMockContext()
	ctx.AddMockDeployment(t, "test-deployment")

	ctx.AddContainerToDeployment(t, "test-deployment", v1.Container{Name: "test-container", Image: "example.com/app:v1"})
	ctx.AddInitContainerToDeployment(t, "test-deployment", v1.Container{Name: "test-init-container", Image: "example.com/init:v1"})
	ctx.AddEphemeralContainerToDeployment(t, "test-deployment", v1.EphemeralContainer{
		EphemeralContainerCommon: v1.EphemeralContainerCommon{Name: "test-ephemeral-container", Image: "example.com/ephemeral:v1"},
	})

	deployment, ok := ctx.objects["test-deployment"].(*appsV1.Deployment)
	require.True(t, ok)

	spec := deployment.Spec.Template.Spec
	require.Len(t, spec.Containers, 1)
	assert.Equal(t, "test-container", spec.Containers[0].Name)
	require.Len(t, spec.InitContainers, 1)
	assert.Equal(t, "test-init-container", spec.InitContainers[0].Name)
	require.Len(t, spec.EphemeralContainers, 1)
	assert.Equal(t, "test-ephemeral-container", spec.EphemeralContainers[0].Name)
}
