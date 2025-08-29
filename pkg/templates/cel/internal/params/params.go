package params

// Params defines the configuration parameters for this template.
type Params struct {
	// Check contains a CEL expression for validation logic. Two predefined variables are available: 'object' (the current Kubernetes object being processed) and 'objects' (all objects being linted).
	// +required
	// +noregex
	Check string
}
