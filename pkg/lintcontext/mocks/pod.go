package mocks

import (
	"testing"

	ocsAppsV1 "github.com/openshift/api/apps/v1"
	"github.com/stretchr/testify/require"
	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	pdbV1 "k8s.io/api/policy/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// AddMockDeployment adds a mock Deployment to LintContext
func (l *MockLintContext) AddMockDeployment(t *testing.T, name string) {
	require.NotEmpty(t, name)
	deployment := &appsV1.Deployment{
		ObjectMeta: metaV1.ObjectMeta{Name: name},
	}
	deployment.SetGroupVersionKind(schema.GroupVersionKind{Group: appsV1.GroupName, Kind: "Deployment"})
	l.objects[name] = deployment
}

// ModifyDeployment modifies a given deployment in the context via the passed function.
func (l *MockLintContext) ModifyDeployment(t *testing.T, name string, f func(deployment *appsV1.Deployment)) {
	dep, ok := l.objects[name].(*appsV1.Deployment)
	require.True(t, ok)
	f(dep)
}

// AddMockDaemonSet adds a mock DaemonSet to LintContext
func (l *MockLintContext) AddMockDaemonSet(t *testing.T, name string) {
	require.NotEmpty(t, name)
	l.objects[name] = &appsV1.DaemonSet{
		ObjectMeta: metaV1.ObjectMeta{Name: name},
	}
}

// ModifyDaemonSet modifies a given DaemonSet in the context via the passed function.
func (l *MockLintContext) ModifyDaemonSet(t *testing.T, name string, f func(ds *appsV1.DaemonSet)) {
	dep, ok := l.objects[name].(*appsV1.DaemonSet)
	require.True(t, ok)
	f(dep)
}

// AddMockDeploymentConfig adds a mock DeploymentConfig to LintContext
func (l *MockLintContext) AddMockDeploymentConfig(t *testing.T, name string) {
	require.NotEmpty(t, name)
	l.objects[name] = &ocsAppsV1.DeploymentConfig{
		ObjectMeta: metaV1.ObjectMeta{Name: name},
	}
}

// ModifyDeploymentConfig modifies a given DeploymentConfig in the context via the passed function.
func (l *MockLintContext) ModifyDeploymentConfig(t *testing.T, name string, f func(ds *ocsAppsV1.DeploymentConfig)) {
	dep, ok := l.objects[name].(*ocsAppsV1.DeploymentConfig)
	require.True(t, ok)
	f(dep)
}

func (l *MockLintContext) AddMockPodDisruptionBudget(t *testing.T, name string) {
	require.NotEmpty(t, name)
	l.objects[name] = &pdbV1.PodDisruptionBudget{
		ObjectMeta: metaV1.ObjectMeta{Name: name},
	}
}

func (l *MockLintContext) ModifyPodDisruptionBudget(t *testing.T, name string, f func(pdb *pdbV1.PodDisruptionBudget)) {
	p, ok := l.objects[name].(*pdbV1.PodDisruptionBudget)
	require.True(t, ok)
	f(p)
}

// AddSecurityContextToDeployment adds a security context to the deployment specified by name
func (l *MockLintContext) AddSecurityContextToDeployment(t *testing.T, deploymentName string,
	securityContext *v1.PodSecurityContext) {
	deployment, ok := l.objects[deploymentName].(*appsV1.Deployment)
	require.True(t, ok, "deployment with name %s not added", deploymentName)
	deployment.Spec.Template.Spec.SecurityContext = securityContext
}

// AddMockReplicationController adds a mock ReplicationController to LintContext
func (l *MockLintContext) AddMockReplicationController(t *testing.T, name string) {
	require.NotEmpty(t, name)
	l.objects[name] = &v1.ReplicationController{
		ObjectMeta: metaV1.ObjectMeta{Name: name},
	}
}

// ModifyReplicationController modifies a given replication controller in the context via the passed function.
func (l *MockLintContext) ModifyReplicationController(t *testing.T, name string, f func(deployment *v1.ReplicationController)) {
	dep, ok := l.objects[name].(*v1.ReplicationController)
	require.True(t, ok)
	f(dep)
}
