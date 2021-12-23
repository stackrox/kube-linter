package params

// Params represents the params accepted by this template.
type Params struct {

	// Key of the forbidden label.
	// +required
	Key string

	// Value of the forbidden label.
	Value string
}
