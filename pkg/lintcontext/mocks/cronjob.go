package mocks

import (
	"testing"

	"github.com/stretchr/testify/require"
	batchV1 "k8s.io/api/batch/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AddMockCronJob adds a mock CronJob to LintContext
func (l *MockLintContext) AddMockCronJob(t *testing.T, name string) {
	require.NotEmpty(t, name)
	l.objects[name] = &batchV1.CronJob{
		ObjectMeta: metaV1.ObjectMeta{Name: name},
	}
}

// ModifyCronJob modifies a given CronJob in the context via the passed function
func (l *MockLintContext) ModifyCronJob(t *testing.T, name string, f func(cronjob *batchV1.CronJob)) {
	dep, ok := l.objects[name].(*batchV1.CronJob)
	require.True(t, ok)
	f(dep)
}
