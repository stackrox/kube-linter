package params

// Params represents the params accepted by this template.
type Params struct {

	// The name of the environment variable.
	// +required
	Name string

	// The value of the environment variable.
	Value string
}
