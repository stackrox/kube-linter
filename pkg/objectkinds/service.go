package objectkinds

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// Service represents Kubernetes Service objects.
	Service = "Service"
)

var (
	serviceGVK = v1.SchemeGroupVersion.WithKind("Service")
)

func init() {
	RegisterObjectKind(Service, MatcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == serviceGVK
	}))
}
