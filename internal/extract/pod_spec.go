package extract

import (
	"reflect"

	"golang.stackrox.io/kube-linter/internal/k8sutil"
	batchV1Beta1 "k8s.io/api/batch/v1beta1"
	coreV1 "k8s.io/api/core/v1"
)

// PodSpec extracts a pod spec from the given object, if available.
func PodSpec(obj k8sutil.Object) (coreV1.PodSpec, bool) {
	switch obj := obj.(type) {
	case *coreV1.Pod:
		return obj.Spec, true
	case *batchV1Beta1.CronJob:
		return obj.Spec.JobTemplate.Spec.Template.Spec, true
	default:
		objValue := reflect.Indirect(reflect.ValueOf(obj))
		spec := objValue.FieldByName("Spec")
		if !spec.IsValid() {
			return coreV1.PodSpec{}, false
		}
		template := spec.FieldByName("Template")
		if !template.IsValid() {
			return coreV1.PodSpec{}, false
		}
		if template.Kind() == reflect.Ptr && !template.IsNil() {
			template = template.Elem()
		}
		podTemplate, ok := template.Interface().(coreV1.PodTemplateSpec)
		if ok {
			return podTemplate.Spec, true
		}
		return coreV1.PodSpec{}, false
	}
}
