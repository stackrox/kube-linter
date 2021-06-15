package extract

import (
	"reflect"

	"golang.stackrox.io/kube-linter/pkg/k8sutil"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// UpdateStrategy will extract the UpdateStrategy struct from a provided
// object if it exists
func UpdateStrategy(obj k8sutil.Object) (interface{}, bool) {
	objValue := reflect.Indirect(reflect.ValueOf(obj))
	spec := objValue.FieldByName("Spec")
	if !spec.IsValid() {
		return nil, false
	}
	strategy := spec.FieldByName("Strategy")
	if !strategy.IsValid() {
		strategy = spec.FieldByName("UpdateStrategy")
		if !strategy.IsValid() {
			return nil, false
		}
	}
	return strategy.Interface(), true
}

// TypeFromUpdateStrategy will extract the Type from a provided
// UpdateStrategy struct if it exists
func TypeFromUpdateStrategy(strategy interface{}) (string, bool) {
	obj := reflect.Indirect(reflect.ValueOf(strategy))
	strategyType := obj.FieldByName("Type")
	if !strategyType.IsValid() {
		return "", false
	}
	return strategyType.String(), true
}

// RollingUpdateFromUpdateStrategy will extract the RollingUpdate struct from a provided
// RollingUpdate struct if it exists
func RollingUpdateFromUpdateStrategy(strategy interface{}) (interface{}, bool) {
	obj := reflect.Indirect(reflect.ValueOf(strategy))
	rollingUpdate := obj.FieldByName("RollingUpdate")
	if !rollingUpdate.IsValid() {
		rollingUpdate = obj.FieldByName("RollingParams")
		if !rollingUpdate.IsValid() {
			return nil, false
		}
	}
	if rollingUpdate.Kind() == reflect.Ptr && !rollingUpdate.IsNil() {
		rollingUpdate = rollingUpdate.Elem()
	}
	return rollingUpdate.Interface(), true
}

// MaxUnavailableFromRollingUpdate will extract the MaxUnavailable field from a provided
// RollingUpdate struct if it exists
func MaxUnavailableFromRollingUpdate(rollingUpdate interface{}) (*intstr.IntOrString, bool) {
	obj := reflect.Indirect(reflect.ValueOf(rollingUpdate))
	if !obj.IsValid() {
		return nil, false
	}
	maxUnavailable := obj.FieldByName("MaxUnavailable")
	if !maxUnavailable.IsValid() {
		return nil, false
	}
	maxUnavailableVal, ok := maxUnavailable.Interface().(*intstr.IntOrString)
	if !ok {
		return nil, false
	}
	return maxUnavailableVal, true
}

// MaxSurgeFromRollingUpdate will extract the MaxSurge field from a provided
// RollingUpdate struct if it exists
func MaxSurgeFromRollingUpdate(rollingUpdate interface{}) (*intstr.IntOrString, bool) {
	obj := reflect.Indirect(reflect.ValueOf(rollingUpdate))
	if !obj.IsValid() {
		return nil, false
	}
	maxSurge := obj.FieldByName("MaxSurge")
	if !maxSurge.IsValid() {
		return nil, false
	}
	maxSurgeVal, ok := maxSurge.Interface().(*intstr.IntOrString)
	if !ok {
		return nil, false
	}
	return maxSurgeVal, true
}
