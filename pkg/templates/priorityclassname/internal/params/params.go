package params

// Params represents the params accepted by this template.
type Params struct {
	// Array of all priority class names that are accepted.
	// +noregex
	// +notnegatable
	AcceptedPriorityClassNames []string `json:"acceptedPriorityClassNames"`
}
