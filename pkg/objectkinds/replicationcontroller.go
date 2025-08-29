package objectkinds

import (
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// ReplicationController represents Kubernetes ReplicationController objects.
	ReplicationController = "ReplicationController"
)

var (
	replicationControllerGVK = coreV1.SchemeGroupVersion.WithKind("ReplicationController")
)

func init() {
	RegisterObjectKind(ReplicationController, MatcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == replicationControllerGVK
	}))
}
