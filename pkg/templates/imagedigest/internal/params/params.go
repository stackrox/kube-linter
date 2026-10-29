package params

// Params represents the params accepted by this template.
type Params struct {

	// list of regular expressions specifying pattern(s) for container images that are exempt from the digest requirement.
	AllowList []string
}
