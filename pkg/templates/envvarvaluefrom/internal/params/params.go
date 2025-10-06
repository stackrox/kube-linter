package params

// Params represents the params accepted by this template.
type Params struct {
	// ignored list => these resources already exist in the cluster.

	// list of regular expressions specifying pattern(s) for secrets that will be ignored.
	IgnoredSecrets []string

	// list of regular expressions specifying pattern(s) for secrets that will be ignored.
	IgnoredConfigMaps []string
}
