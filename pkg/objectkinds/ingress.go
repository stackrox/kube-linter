package objectkinds

import (
	v1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// Ingress represents Kubernetes Ingress objects.
	Ingress = "Ingress"
)

var (
	ingressGVK = v1.SchemeGroupVersion.WithKind(Ingress)
)

func init() {
	RegisterObjectKind(Ingress, MatcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == ingressGVK
	}))
}

// GetIngressAPIVersion returns Ingress's apiversion
func GetIngressAPIVersion() string {
	return ingressGVK.GroupVersion().String()
}
