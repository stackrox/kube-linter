/*
Copyright 2021 The KEDA Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConditionType specifies the available conditions for the resource
type ConditionType string

const (
	// ConditionReady specifies that the resource is ready.
	// For long-running resources.
	ConditionReady ConditionType = "Ready"
	// ConditionActive specifies that the resource has finished.
	// For resource which run to completion.
	ConditionActive ConditionType = "Active"
	// ConditionFallback specifies that the resource has a fallback active.
	ConditionFallback ConditionType = "Fallback"
	// ConditionPaused specifies that the resource is paused.
	ConditionPaused ConditionType = "Paused"
)

const (
	// ScaledObjectConditionReadySuccessReason defines the default Reason for correct ScaledObject
	ScaledObjectConditionReadySuccessReason = "ScaledObjectReady"
	// ScaledObjectConditionReadySuccessMessage defines the default Message for correct ScaledObject
	ScaledObjectConditionReadySuccessMessage = "ScaledObject is defined correctly and is ready for scaling"
	// ScaledObjectConditionPausedReason defines the default Reason for paused ScaledObject
	ScaledObjectConditionPausedReason = "ScaledObjectPaused"
	// ScaledObjectConditionPausedMessage defines the default Message for paused ScaledObject
	ScaledObjectConditionPausedMessage = "ScaledObject is paused"
)

const (
	// ScaledJobConditionPausedReason defines the default Reason for paused ScaledJob
	ScaledJobConditionPausedReason = "ScaledJobPaused"
	// ScaledJobConditionPausedReason defines the default Reason for paused ScaledJob
	ScaledJobConditionUnpausedReason = "ScaledJobUnpaused"
	// ScaledJobConditionPausedMessage defines the default Message for paused ScaledJob
	ScaledJobConditionPausedMessage = "ScaledJob is paused"
	// ScaledJobConditionPausedMessage defines the default Message for paused ScaledJob
	ScaledJobConditionUnpausedMessage = "ScaledJob is unpaused"
)

// Condition to store the condition state
type Condition struct {
	// Type of condition
	// +required
	Type ConditionType `json:"type" description:"type of status condition"`

	// Status of the condition, one of True, False, Unknown.
	// +required
	Status metav1.ConditionStatus `json:"status" description:"status of the condition, one of True, False, Unknown"`

	// The reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty" description:"one-word CamelCase reason for the condition's last transition"`

	// A human readable message indicating details about the transition.
	// +optional
	Message string `json:"message,omitempty" description:"human-readable message indicating details about last transition"`
}

// Conditions an array representation to store multiple Conditions
type Conditions []Condition

// IsTrue is true if the condition is True
func (c *Condition) IsTrue() bool {
	if c == nil {
		return false
	}
	return c.Status == metav1.ConditionTrue
}

// IsFalse is true if the condition is False
func (c *Condition) IsFalse() bool {
	if c == nil {
		return false
	}
	return c.Status == metav1.ConditionFalse
}

// IsUnknown is true if the condition is Unknown
func (c *Condition) IsUnknown() bool {
	if c == nil {
		return true
	}
	return c.Status == metav1.ConditionUnknown
}
