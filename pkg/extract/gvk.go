package extract

import (
	"golang.stackrox.io/kube-linter/pkg/k8sutil"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GVK extracts the GroupVersionKind of an object.
func GVK(object k8sutil.Object) schema.GroupVersionKind {
	return object.GetObjectKind().GroupVersionKind()
}
