package params

// Params represents the params accepted by this template.
type Params struct {

	// A regex specifying the required service account to match.
	// +required
	ServiceAccount string
}
