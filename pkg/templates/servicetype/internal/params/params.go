package params

// Params represents the params accepted by this template.
type Params struct {
	// An array of service types that should not be used
	// +noregex
	// +notnegatable
	ForbiddenServiceTypes []string `json:"forbiddenServiceTypes"`
}
