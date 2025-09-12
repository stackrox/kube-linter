package objectkinds

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	PersistentVolumeClaim = "PersistentVolumeClaim"
)

var (
	persistentvolumeclaimGVK = v1.SchemeGroupVersion.WithKind("PersistentVolumeClaim")
)

func init() {
	RegisterObjectKind(PersistentVolumeClaim, MatcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == persistentvolumeclaimGVK
	}))
}

func GetPersistentVolumeClaimAPIVersion() string {
	return persistentvolumeclaimGVK.GroupVersion().String()
}
