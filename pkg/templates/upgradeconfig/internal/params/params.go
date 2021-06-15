package params

// Params represents the params accepted by this template.
type Params struct {
	StrategyTypeRegex  string
	MaxPodsUnavailable string
	MinPodsUnavailable string
	MaxSurge           string
	MinSurge           string
}
