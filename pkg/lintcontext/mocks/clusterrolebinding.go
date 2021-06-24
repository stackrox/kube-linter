package mocks

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	rbacV1 "k8s.io/api/rbac/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AddMockClusterRoleBinding adds a mock ClusterRoleBinding to LintContext
func (l *MockLintContext) AddMockClusterRoleBinding(t *testing.T, name string) {
	require.NotEmpty(t, name)
	l.objects[name] = &rbacV1.ClusterRoleBinding{
		TypeMeta: metaV1.TypeMeta{
			Kind:       objectkinds.ClusterRoleBinding,
			APIVersion: objectkinds.GetClusterRoleBindingAPIVersion(),
		},
		ObjectMeta: metaV1.ObjectMeta{Name: name},
		Subjects:   []rbacV1.Subject{},
		RoleRef:    rbacV1.RoleRef{},
	}
}

// ModifyClusterRoleBinding modifies a given ClusterRoleBinding in the context via the passed function.
func (l *MockLintContext) ModifyClusterRoleBinding(t *testing.T, name string, f func(clusterrolebinding *rbacV1.ClusterRoleBinding)) {
	crb, ok := l.objects[name].(*rbacV1.ClusterRoleBinding)
	require.True(t, ok)
	f(crb)
}
