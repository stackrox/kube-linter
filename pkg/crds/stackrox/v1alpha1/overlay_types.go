// Copyright 2019 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

// K8sObjectOverlay is an overlay that applies a set of patches to a resource.
// It targets a resource by its API version, kind, and name, and applies
// a list of patches to this resource.
//
// # Examples
//
// ## Adding an annotation to a resource
//
//	apiVersion: v1
//	kind: ServiceAccount
//	name: central
//	patches:
//	- path: metadata.annotations.eks\.amazonaws\.com/role-arn
//	  value: "\"arn:aws:iam:1234:role\""
//
// ## Adding an environment variable to a deployment
//
//	apiVersion: apps/v1
//	kind: Deployment
//	name: central
//	patches:
//	- path: spec.template.spec.containers[name:central].env[-1]
//	  value: |
//	    name: MY_ENV_VAR
//	    value: value
//
// ## Adding an ingress to a network policy
//
//	apiVersion: networking.k8s.io/v1
//	kind: NetworkPolicy
//	name: allow-ext-to-central
//	patches:
//	- path: spec.ingress[-1]
//	  value: |
//	    ports:
//	    - port: 999
//	      protocol: TCP
//
// ## Changing the value of a ConfigMap
//
//	apiVersion: v1
//	kind: ConfigMap
//	name: central-endpoints
//	patches:
//	- path: data.endpoints\.yaml:
//	  verbatim: |
//	    disableDefault: false
//
// ## Adding a container to a deployment
//
//	apiVersion: apps/v1
//	kind: Deployment
//	name: central
//	patches:
//	  - path: spec.template.spec.containers[-1]
//	    value: |
//	      name: nginx
//	      image: nginx
//	      ports:
//	      - containerPort: 8000
//	        name: http
//	        protocol: TCP
type K8sObjectOverlay struct {
	// Resource API version.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="API Version",order=1
	APIVersion string `json:"apiVersion,omitempty"`
	// Resource kind.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Kind",order=2
	Kind string `json:"kind,omitempty"`
	// Name of resource.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Name",order=3
	Name string `json:"name,omitempty"`
	// Optional marks the overlay as optional.
	// When Optional is true, and the specified resource does not exist in the output manifests, the overlay will be skipped, and a warning will be logged.
	// When Optional is false, and the specified resource does not exist in the output manifests, an error will be thrown.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Optional",order=4
	Optional bool `json:"optional,omitempty"`
	// List of patches to apply to resource.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Patches",order=5
	Patches []*K8sObjectOverlayPatch `json:"patches,omitempty"`
}

// K8sObjectOverlayPatch defines a patch to apply to a resource.
type K8sObjectOverlayPatch struct {
	// Path of the form a.[key1:value1].b.[:value2]
	// Where [key1:value1] is a selector for a key-value pair to identify a list element and [:value] is a value
	// selector to identify a list element in a leaf list.
	// All path intermediate nodes must exist.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Path",order=1
	Path string `json:"path,omitempty"`
	// Value to add, delete or replace.
	// For add, the path should be a new leaf.
	// For delete, value should be unset.
	// For replace, path should reference an existing node.
	// All values are strings but are converted into appropriate type based on schema.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Value",order=2
	Value string `json:"value,omitempty"`
	// Verbatim value to add, delete or replace.
	// Same as Value, but the content is not interpreted as YAML and is treated as a literal string instead.
	// At least one of Value and Verbatim must be empty.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Verbatim",order=3
	Verbatim string `json:"verbatim,omitempty"`
}
