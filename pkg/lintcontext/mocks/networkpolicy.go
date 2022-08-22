package mocks

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	networkingV1 "k8s.io/api/networking/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AddMockNetworkPolicy  adds a mock NetworkPolicy to LintContext
func (l *MockLintContext) AddMockNetworkPolicy(t *testing.T, name string) {
	require.NotEmpty(t, name)
	l.objects[name] = &networkingV1.NetworkPolicy{
		TypeMeta: metaV1.TypeMeta{
			Kind:       objectkinds.NetworkPolicy,
			APIVersion: objectkinds.GetNetworkPolicyAPIVersion(),
		},
		ObjectMeta: metaV1.ObjectMeta{Name: name},
		Spec:       networkingV1.NetworkPolicySpec{},
	}
}

// ModifyNetworkPolicy modifies a given networkpolicy in the context via the passed function.
func (l *MockLintContext) ModifyNetworkPolicy(t *testing.T, name string, f func(networkpolicy *networkingV1.NetworkPolicy)) {
	r, ok := l.objects[name].(*networkingV1.NetworkPolicy)
	require.True(t, ok)
	f(r)
}
