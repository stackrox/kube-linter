package params

// Params represents the params accepted by this template.
type Params struct {

	// Key of the forbidden annotation.
	// +required
	Key string

	// Value of the forbidden annotation.
	Value string
}
