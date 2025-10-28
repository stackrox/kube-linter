package params

// Params defines the configuration parameters for this template.
type Params struct {
	// SchemaLocations contains locations of schemas to use. See: https://github.com/yannh/kubeconform/tree/master?tab=readme-ov-file#overriding-schemas-location
	// +noregex
	SchemaLocations []string
	// Cache specifies the folder to cache schemas downloaded via HTTP.
	// +noregex
	Cache string
	// SkipKinds lists resource kinds to ignore during validation.
	// +noregex
	SkipKinds []string
	// RejectKinds lists resource kinds to reject during validation.
	// +noregex
	RejectKinds []string
	// KubernetesVersion specifies the Kubernetes version - must match one in https://github.com/instrumenta/kubernetes-json-schema
	// +noregex
	KubernetesVersion string
	// Strict enables strict validation that will error if resources contain undocumented fields.
	Strict bool
	// IgnoreMissingSchemas will skip validation for resources if no schema can be found.
	IgnoreMissingSchemas bool
}
