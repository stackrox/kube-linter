package objectkinds

import (
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// Pod represents Kubernetes Pod objects.
	Pod = "Pod"
)

var (
	podGVK = coreV1.SchemeGroupVersion.WithKind("Pod")
)

func init() {
	RegisterObjectKind(Pod, MatcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == podGVK
	}))
}
