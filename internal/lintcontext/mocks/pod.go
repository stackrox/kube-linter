package mocks

import (
	"testing"

	"github.com/stretchr/testify/require"
	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AddMockDeployment adds a mock Deployment to LintContext
func (l *MockLintContext) AddMockDeployment(t *testing.T, name string) {
	require.NotEmpty(t, name)
	l.objects[name] = &appsV1.Deployment{
		ObjectMeta: metaV1.ObjectMeta{Name: name},
	}
}

// ModifyDeployment modifies a given deployment in the context via the passed function.
func (l *MockLintContext) ModifyDeployment(t *testing.T, name string, f func(deployment *appsV1.Deployment)) {
	dep, ok := l.objects[name].(*appsV1.Deployment)
	require.True(t, ok)
	f(dep)
}

// AddSecurityContextToDeployment adds a security context to the deployment specified by name
func (l *MockLintContext) AddSecurityContextToDeployment(t *testing.T, deploymentName string,
	securityContext *v1.PodSecurityContext) {
	deployment, ok := l.objects[deploymentName].(*appsV1.Deployment)
	require.True(t, ok, "deployment with name %s not added", deploymentName)
	deployment.Spec.Template.Spec.SecurityContext = securityContext
}
