package params

// Params represents the params accepted by this template.
type Params struct {
	// A regular expression the defines the type of update
	// strategy allowed.
	// +required
	StrategyTypeRegex string

	// The maximum value that be set in a RollingUpdate
	// configuration for the MaxUnavailable.  This can be
	// an integer or a percent.
	// +optional
	MaxPodsUnavailable string

	// The minimum value that be set in a RollingUpdate
	// configuration for the MaxUnavailable.  This can be
	// an integer or a percent.
	// +optional
	MinPodsUnavailable string

	// The maximum value that be set in a RollingUpdate
	// configuration for the MaxSurge.  This can be
	// an integer or a percent.
	// +optional
	MaxSurge string

	// The minimum value that be set in a RollingUpdate
	// configuration for the MaxSurge.  This can be
	// an integer or a percent.
	// +optional
	MinSurge string
}
