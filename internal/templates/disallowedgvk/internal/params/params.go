package params

// Params represents the params accepted by this template.
type Params struct {

	// The disallowed object group.
	// +example=apps
	Group string `json:"group"`

	// The disallowed object API version.
	// +example=v1
	// +example=v1beta1
	Version string

	// The disallowed kind.
	// +example=Deployment
	// +example=DaemonSet
	Kind string
}
