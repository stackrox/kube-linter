package mocks

import (
	"testing"

	"github.com/stretchr/testify/require"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AddMockService adds a mock Service to LintContext
func (l *MockLintContext) AddMockService(t *testing.T, name string) {
	require.NotEmpty(t, name)
	l.objects[name] = &coreV1.Service{
		ObjectMeta: metaV1.ObjectMeta{Name: name},
	}
}

// ModifyService modifies a given service in the context via the passed function
func (l *MockLintContext) ModifyService(t *testing.T, name string, f func(service *coreV1.Service)) {
	dep, ok := l.objects[name].(*coreV1.Service)
	require.True(t, ok)
	f(dep)
}
