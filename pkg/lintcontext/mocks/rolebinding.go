package mocks

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	rbacV1 "k8s.io/api/rbac/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AddMockRoleBinding adds a mock RoleBinding to LintContext
func (l *MockLintContext) AddMockRoleBinding(t *testing.T, name, namespace string) {
	require.NotEmpty(t, name)
	l.objects[name] = &rbacV1.RoleBinding{
		TypeMeta: metaV1.TypeMeta{
			Kind:       objectkinds.RoleBinding,
			APIVersion: objectkinds.GetRoleBindingAPIVersion(),
		},
		ObjectMeta: metaV1.ObjectMeta{Name: name, Namespace: namespace},
		Subjects:   []rbacV1.Subject{},
		RoleRef:    rbacV1.RoleRef{},
	}
}

// ModifyRoleBinding modifies a given RoleBinding in the context via the passed function.
func (l *MockLintContext) ModifyRoleBinding(t *testing.T, name string, f func(rolebinding *rbacV1.RoleBinding)) {
	rb, ok := l.objects[name].(*rbacV1.RoleBinding)
	require.True(t, ok)
	f(rb)
}
