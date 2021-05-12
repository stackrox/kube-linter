package extract

import (
	"reflect"

	"golang.stackrox.io/kube-linter/pkg/k8sutil"
	batchV1Beta1 "k8s.io/api/batch/v1beta1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodTemplateSpec extracts a pod template spec from the given object, if available.
func PodTemplateSpec(obj k8sutil.Object) (coreV1.PodTemplateSpec, bool) {
	switch obj := obj.(type) {
	case *coreV1.Pod:
		return coreV1.PodTemplateSpec{
			ObjectMeta: obj.ObjectMeta,
			Spec:       obj.Spec,
		}, true
	case *batchV1Beta1.CronJob:
		return obj.Spec.JobTemplate.Spec.Template, true
	default:
		objValue := reflect.Indirect(reflect.ValueOf(obj))
		spec := objValue.FieldByName("Spec")
		if !spec.IsValid() {
			return coreV1.PodTemplateSpec{}, false
		}
		template := spec.FieldByName("Template")
		if !template.IsValid() {
			return coreV1.PodTemplateSpec{}, false
		}
		if template.Kind() == reflect.Ptr && !template.IsNil() {
			template = template.Elem()
		}
		podTemplate, ok := template.Interface().(coreV1.PodTemplateSpec)
		if ok {
			return podTemplate, true
		}
		return coreV1.PodTemplateSpec{}, false
	}
}

// PodSpec extracts a pod spec from the given object, if available.
func PodSpec(obj k8sutil.Object) (coreV1.PodSpec, bool) {
	podTemplateSpec, found := PodTemplateSpec(obj)
	if !found {
		return coreV1.PodSpec{}, false
	}
	return podTemplateSpec.Spec, true
}

// Selector extracts a selector from the given object, if available.
func Selector(obj k8sutil.Object) (*metaV1.LabelSelector, bool) {
	switch obj := obj.(type) {
	case *batchV1Beta1.CronJob:
		selector := obj.Spec.JobTemplate.Spec.Selector
		return selector, true
	default:
		objValue := reflect.Indirect(reflect.ValueOf(obj))
		spec := objValue.FieldByName("Spec")
		if !spec.IsValid() {
			return nil, false
		}
		selector := spec.FieldByName("Selector")
		if !selector.IsValid() {
			return nil, false
		}
		labelSelector, ok := selector.Interface().(*metaV1.LabelSelector)
		if ok {
			return labelSelector, true
		}
	}
	return nil, false
}

// Replicas extracts replicas from the given object, if available.
func Replicas(obj k8sutil.Object) (int32, bool) {
	objValue := reflect.Indirect(reflect.ValueOf(obj))
	spec := objValue.FieldByName("Spec")
	if !spec.IsValid() {
		return 0, false
	}
	replicas := spec.FieldByName("Replicas")
	if !replicas.IsValid() {
		return 0, false
	}
	numReplicas, ok := replicas.Interface().(*int32)
	if ok {
		if numReplicas != nil {
			return *numReplicas, true
		}
		// If numReplicas is a `nil` pointer, then it defaults to 1.
		return 1, true
	}
	return 0, false
}
