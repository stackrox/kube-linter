package k8sutil

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Object is a combination of `runtime.Object` and `metav1.Object`.
type Object interface {
	runtime.Object
	metav1.Object
}
