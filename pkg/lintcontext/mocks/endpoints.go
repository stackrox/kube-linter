package mocks

import (
	"testing"

	"github.com/stretchr/testify/require"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AddMockEndpoints adds a mock Endpoint to LintContext
func (l *MockLintContext) AddMockEndpoints(t *testing.T, name string) {
	require.NotEmpty(t, name)
	l.objects[name] = &coreV1.Endpoints{
		ObjectMeta: metaV1.ObjectMeta{Name: name},
	}
}
