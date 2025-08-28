package params

// Params represents the params accepted by this template.
type Params struct {
	// Check is a CEL expression used to validate the subject and objects.
	// +required
	// +noregex
	Check string
}
