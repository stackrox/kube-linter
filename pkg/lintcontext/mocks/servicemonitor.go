package mocks

import (
	"testing"

	k8sMonitoring "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/stretchr/testify/require"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AddMockServiceMonitor adds a mock ServiceMonitor to LintContext
func (l *MockLintContext) AddMockServiceMonitor(t *testing.T, name string) {
	require.NotEmpty(t, name)
	l.objects[name] = &k8sMonitoring.ServiceMonitor{
		TypeMeta: metaV1.TypeMeta{
			Kind:       objectkinds.ServiceMonitor,
			APIVersion: objectkinds.GetServiceMonitorAPIVersion(),
		},
		ObjectMeta: metaV1.ObjectMeta{Name: name},
		Spec:       k8sMonitoring.ServiceMonitorSpec{},
	}
}

// ModifyServiceMonitor modifies a given servicemonitor in the context via the passed function
func (l *MockLintContext) ModifyServiceMonitor(t *testing.T, name string, f func(servicemonitor *k8sMonitoring.ServiceMonitor)) {
	r, ok := l.objects[name].(*k8sMonitoring.ServiceMonitor)
	require.True(t, ok)
	f(r)
}
