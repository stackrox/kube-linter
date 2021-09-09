package params

// Params represents the params accepted by this template.
type Params struct {

	// The minimum value of "MaxUnavailable" field in the PodDisruptionBudget object
	MinimumMaxUnavailableCriteria string
}
