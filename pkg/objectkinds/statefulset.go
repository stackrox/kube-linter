package objectkinds

import (
	appsV1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// StatefulSet represents Kubernetes StatefulSet objects.
	StatefulSet = "StatefulSet"
)

var (
	statefulSetGVK = appsV1.SchemeGroupVersion.WithKind("StatefulSet")
)

func init() {
	RegisterObjectKind(StatefulSet, MatcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == statefulSetGVK
	}))
}
