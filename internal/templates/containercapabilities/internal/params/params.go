package params

// Params represents the params accepted by this template.
type Params struct {

	// List of capabilities that needs to be removed from containers.
	// +noregex
	// +notnegatable
	ForbiddenCapabilities []string `json:"forbiddenCapabilities"`

	// List of capabilities that are exceptions to the above list. This should only be filled
	// when the above contains "all", and is used to forgive capabilities in ADD list.
	// +noregex
	// +notnegatable
	Exceptions []string `json:"exceptions"`
}
