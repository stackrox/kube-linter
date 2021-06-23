package extract

import (
	"reflect"

	"golang.stackrox.io/kube-linter/pkg/k8sutil"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// UpdateStrategyValues contains testable data from an UpdateStrategy struct
type UpdateStrategyValues struct {
	Type                 string
	TypeExists           bool
	RollingConfigExists  bool
	RollingConfigValid   bool
	MaxUnavailableExists bool
	MaxUnavailable       *intstr.IntOrString
	MaxSurgeExists       bool
	MaxSurge             *intstr.IntOrString
}

// UpdateStrategy will extract the data from an UpdateStrategy into a common struct
func UpdateStrategy(obj k8sutil.Object) (*UpdateStrategyValues, bool) {
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
	strategyType, typeFound := typeFromUpdateStrategy(strategy)
	rollingUpdate, rollingUpdateFound := rollingUpdateFromUpdateStrategy(strategy)
	maxUnavailable, maxUnavailableFound := maxUnavailableFromRollingUpdate(rollingUpdate)
	maxSurge, maxSurgeFound := maxSurgeFromRollingUpdate(rollingUpdate)

	return &UpdateStrategyValues{
		Type:                 strategyType,
		TypeExists:           typeFound,
		RollingConfigExists:  rollingUpdateFound,
		RollingConfigValid:   reflect.Indirect(rollingUpdate).IsValid(),
		MaxUnavailable:       maxUnavailable,
		MaxUnavailableExists: maxUnavailableFound,
		MaxSurge:             maxSurge,
		MaxSurgeExists:       maxSurgeFound,
	}, true
}

// typeFromUpdateStrategy will extract the Type from a provided
// UpdateStrategy struct if it exists
func typeFromUpdateStrategy(strategy reflect.Value) (string, bool) {
	obj := reflect.Indirect(strategy)
	strategyType := obj.FieldByName("Type")
	if !strategyType.IsValid() {
		return "", false
	}
	return strategyType.String(), true
}

// rollingUpdateFromUpdateStrategy will extract the RollingUpdate struct from a provided
// RollingUpdate struct if it exists
func rollingUpdateFromUpdateStrategy(strategy reflect.Value) (reflect.Value, bool) {
	obj := reflect.Indirect(strategy)
	rollingUpdate := obj.FieldByName("RollingUpdate")
	if !rollingUpdate.IsValid() {
		rollingUpdate = obj.FieldByName("RollingParams")
		if !rollingUpdate.IsValid() {
			return rollingUpdate, false
		}
	}
	if rollingUpdate.Kind() == reflect.Ptr && !rollingUpdate.IsNil() {
		rollingUpdate = rollingUpdate.Elem()
	}
	return rollingUpdate, true
}

// maxUnavailableFromRollingUpdate will extract the MaxUnavailable field from a provided
// RollingUpdate struct if it exists
func maxUnavailableFromRollingUpdate(rollingUpdate reflect.Value) (*intstr.IntOrString, bool) {
	obj := reflect.Indirect(rollingUpdate)
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

// maxSurgeFromRollingUpdate will extract the MaxSurge field from a provided
// RollingUpdate struct if it exists
func maxSurgeFromRollingUpdate(rollingUpdate reflect.Value) (*intstr.IntOrString, bool) {
	obj := reflect.Indirect(rollingUpdate)
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
