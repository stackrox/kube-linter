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
	ingresssGVK = v1.SchemeGroupVersion.WithKind(Ingress)
)

func init() {
	registerObjectKind(Ingress, matcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == ingresssGVK
	}))
}

// GetIngressAPIVersion returns Ingress's apiversion
func GetIngressAPIVersion() string {
	return ingresssGVK.GroupVersion().String()
}
