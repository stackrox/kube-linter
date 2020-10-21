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
	// a number of MB.
	LowerBoundMB int `json:"lowerBoundMB"`

	// The upper bound of the requirement (inclusive), specified as
	// a number of MB.
	// If not specified, it is treated as "no upper bound".
	UpperBoundMB *int `json:"upperBoundMB"`
}
