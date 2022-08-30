package objectkinds

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// ServiceAccount represents Kubernetes ServiceAccount objects.
	ServiceAccount = "ServiceAccount"
)

var (
	serviceAccountGVK = v1.SchemeGroupVersion.WithKind("ServiceAccount")
)

func init() {
	RegisterObjectKind(ServiceAccount, MatcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == serviceAccountGVK
	}))
}
