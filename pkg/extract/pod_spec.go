package extract

import (
	"reflect"

	ocsAppsV1 "github.com/openshift/api/apps/v1"
	"golang.stackrox.io/kube-linter/pkg/extract/customtypes"
	"golang.stackrox.io/kube-linter/pkg/k8sutil"
	batchV1 "k8s.io/api/batch/v1"
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
	case *batchV1.CronJob:
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
func PodSpec(obj k8sutil.Object) (customtypes.PodSpec, bool) {
	podTemplateSpec, found := PodTemplateSpec(obj)
	if !found {
		return customtypes.PodSpec{}, false
	}
	return customtypes.PodSpec{PodSpec: podTemplateSpec.Spec}, true
}

// Selector extracts a selector from the given object, if available.
func Selector(obj k8sutil.Object) (*metaV1.LabelSelector, bool) {
	switch obj := obj.(type) {
	case *ocsAppsV1.DeploymentConfig:
		return &metaV1.LabelSelector{MatchLabels: obj.Spec.Selector}, true
	case *batchV1Beta1.CronJob:
		return obj.Spec.JobTemplate.Spec.Selector, true
	case *batchV1.CronJob:
		return obj.Spec.JobTemplate.Spec.Selector, true
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
	// DeploymentConfigs are treated specially because the number of replicas is
	// an int32, not a *int32.
	if depConfig, isDepConfig := obj.(*ocsAppsV1.DeploymentConfig); isDepConfig {
		return depConfig.Spec.Replicas, true
	}
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
