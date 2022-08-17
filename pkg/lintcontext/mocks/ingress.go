package mocks

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	networkingV1 "k8s.io/api/networking/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (l *MockLintContext) AddMockIngress(t *testing.T, name string) {
	require.NotEmpty(t, name)
	l.objects[name] = &networkingV1.Ingress{
		TypeMeta: metaV1.TypeMeta{
			Kind:       objectkinds.Ingress,
			APIVersion: objectkinds.GetIngressAPIVersion(),
		},
		ObjectMeta: metaV1.ObjectMeta{Name: name},
		Spec:       networkingV1.IngressSpec{},
	}
}

// ModifyIngress modifies a given networkpolicy in the context via the passed function.
func (l *MockLintContext) ModifyIngress(t *testing.T, name string, f func(ingress *networkingV1.Ingress)) {
	r, ok := l.objects[name].(*networkingV1.Ingress)
	require.True(t, ok)
	f(r)
}
