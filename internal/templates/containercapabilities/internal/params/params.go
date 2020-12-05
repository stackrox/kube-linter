package params

// Params represents the params accepted by this template.
type Params struct {

	// List of capabilities that are required to be dropped by containers.
	// +notnegatable
	RequiredDrops []string `json:"requiredDrops"`

	// List of capabilities that are forbidden to be added to containers.
	// +notnegatable
	ForbiddenAdds []string `json:"forbiddenAdds"`
}
