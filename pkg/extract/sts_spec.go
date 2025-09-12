package extract

import (
	"reflect"

	"golang.stackrox.io/kube-linter/pkg/k8sutil"
	appsV1 "k8s.io/api/apps/v1"
)

func StatefulSetSpec(obj k8sutil.Object) (appsV1.StatefulSetSpec, bool) {
	if obj == nil {
		return appsV1.StatefulSetSpec{}, false
	}

	switch obj := obj.(type) {
	case *appsV1.StatefulSet:
		return obj.Spec, true
	default:
		kind := obj.GetObjectKind().GroupVersionKind().Kind
		if kind != "StatefulSet" {
			return appsV1.StatefulSetSpec{}, false
		}

		objValue := reflect.Indirect(reflect.ValueOf(obj))
		spec := objValue.FieldByName("Spec")
		if !spec.IsValid() {
			return appsV1.StatefulSetSpec{}, false
		}
		statefulSetSpec, ok := spec.Interface().(appsV1.StatefulSetSpec)
		if ok {
			return statefulSetSpec, true
		}
		return appsV1.StatefulSetSpec{}, false
	}
}
