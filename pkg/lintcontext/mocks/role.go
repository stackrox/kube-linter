package mocks

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	rbacV1 "k8s.io/api/rbac/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AddMockRole adds a mock Role to LintContext
func (l *MockLintContext) AddMockRole(t *testing.T, name, namespace string) {
	require.NotEmpty(t, name)
	l.objects[name] = &rbacV1.Role{
		TypeMeta: metaV1.TypeMeta{
			Kind:       objectkinds.Role,
			APIVersion: objectkinds.GetRoleAPIVersion(),
		},
		ObjectMeta: metaV1.ObjectMeta{Name: name, Namespace: namespace},
		Rules:      []rbacV1.PolicyRule{},
	}
}

// ModifyRole modifies a given Role in the context via the passed function.
func (l *MockLintContext) ModifyRole(t *testing.T, name string, f func(role *rbacV1.Role)) {
	r, ok := l.objects[name].(*rbacV1.Role)
	require.True(t, ok)
	f(r)
}
