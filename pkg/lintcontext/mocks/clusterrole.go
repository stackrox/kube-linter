package mocks

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	rbacV1 "k8s.io/api/rbac/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AddMockClusterRole adds a mock ClusterRole to LintContext
func (l *MockLintContext) AddMockClusterRole(t *testing.T, name string) {
	require.NotEmpty(t, name)
	l.objects[name] = &rbacV1.ClusterRole{
		TypeMeta: metaV1.TypeMeta{
			Kind:       objectkinds.ClusterRole,
			APIVersion: objectkinds.GetClusterRoleAPIVersion(),
		},
		ObjectMeta:      metaV1.ObjectMeta{Name: name},
		Rules:           []rbacV1.PolicyRule{},
		AggregationRule: &rbacV1.AggregationRule{},
	}
}

// ModifyClusterRole modifies a given clusterrole in the context via the passed function.
func (l *MockLintContext) ModifyClusterRole(t *testing.T, name string, f func(clusterrole *rbacV1.ClusterRole)) {
	r, ok := l.objects[name].(*rbacV1.ClusterRole)
	require.True(t, ok)
	f(r)
}
