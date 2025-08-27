package params

// Params represents the params accepted by this template.
type Params struct {

	// Check is a CEL Expression validate the subject and objects
	Check string

	// Filter is a CEL Expression that returns true if GroupVersionKind should be inspected by Check.
	Filter string
}
