package objectkinds

import (
	appsV1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// Deployment represents Kubernetes Deployment objects.
	Deployment = "Deployment"
)

var (
	deploymentGVK = appsV1.SchemeGroupVersion.WithKind("Deployment")
)

func init() {
	RegisterObjectKind(Deployment, MatcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == deploymentGVK
	}))
}
