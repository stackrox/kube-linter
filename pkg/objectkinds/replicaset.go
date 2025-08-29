package objectkinds

import (
	appsV1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// ReplicaSet represents Kubernetes ReplicaSet objects.
	ReplicaSet = "ReplicaSet"
)

var (
	replicaSetGVK = appsV1.SchemeGroupVersion.WithKind("ReplicaSet")
)

func init() {
	RegisterObjectKind(ReplicaSet, MatcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == replicaSetGVK
	}))
}
