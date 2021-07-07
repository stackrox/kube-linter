package params

// Params represents the params accepted by this template.
type Params struct {

	// list of regular expressions for blocked or bad container image tags
	BlockList []string
}
