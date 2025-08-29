package objectkinds

import (
	appsV1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// DaemonSet represents Kubernetes DaemonSet objects.
	DaemonSet = "DaemonSet"
)

var (
	daemonSetGVK = appsV1.SchemeGroupVersion.WithKind("DaemonSet")
)

func init() {
	RegisterObjectKind(DaemonSet, MatcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == daemonSetGVK
	}))
}
