package extract

import (
	"golang.stackrox.io/kube-linter/pkg/k8sutil"
)

// Labels extracts labels from the given object.
func Labels(object k8sutil.Object) map[string]string {
	return object.GetLabels()
}

// Annotations extracts annotations from the given object.
func Annotations(object k8sutil.Object) map[string]string {
	return object.GetAnnotations()
}
