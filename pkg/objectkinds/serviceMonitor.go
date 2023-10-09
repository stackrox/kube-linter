package objectkinds

import (
	k8sMonitoring "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// ServiceMonitor represents Prometheus Service Monitor objects.
	ServiceMonitor = k8sMonitoring.ServiceMonitorsKind
)

var (
	serviceMonitorGVK = k8sMonitoring.SchemeGroupVersion.WithKind(ServiceMonitor)
)

func init() {
	RegisterObjectKind(ServiceMonitor, MatcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == serviceMonitorGVK
	}))
}

// GetServiceMonitorAPIVersion returns servicemonitor's apiversion
func GetServiceMonitorAPIVersion() string {
	return serviceMonitorGVK.GroupVersion().String()
}
