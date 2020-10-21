package params

// Params represents the params accepted by this template.
type Params struct {

	// The type of requirement. Use any to apply to both requests and limits.
	// +enum=request
	// +enum=limit
	// +enum=any
	// +required
	RequirementsType string

	// The lower bound of the requirement (inclusive), specified as
	// a number of milli-cores.
	// If not specified, it is treated as a lower bound of zero.
	LowerBoundMillis int `json:"lowerBoundMillis"`

	// The upper bound of the requirement (inclusive), specified as
	// a number of milli-cores.
	// If not specified, it is treated as "no upper bound".
	UpperBoundMillis *int `json:"upperBoundMillis"`
}
