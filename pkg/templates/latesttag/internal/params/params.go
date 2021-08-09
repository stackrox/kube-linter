package params

// Params represents the params accepted by this template.
type Params struct {

	// list of regular expressions specifying pattern(s) for container images that will be blocked. */
	BlockList []string

	// list of regular expressions specifying pattern(s) for container images that will be allowed.
	AllowList []string
}
