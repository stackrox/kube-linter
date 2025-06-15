package mocks

import (
	"testing"

	"github.com/stretchr/testify/require"
	batchV1 "k8s.io/api/batch/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AddMockJob adds a mock Job to LintContext
func (l *MockLintContext) AddMockJob(t *testing.T, name string) {
	require.NotEmpty(t, name)
	l.objects[name] = &batchV1.Job{
		ObjectMeta: metaV1.ObjectMeta{Name: name},
	}
}

// ModifyJob modifies a given Job in the context via the passed function
func (l *MockLintContext) ModifyJob(t *testing.T, name string, f func(job *batchV1.Job)) {
	dep, ok := l.objects[name].(*batchV1.Job)
	require.True(t, ok)
	f(dep)
}
