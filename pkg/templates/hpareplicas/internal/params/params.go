package params

// Params represents the params accepted by this template.
type Params struct {

	// The minimum number of replicas a HorizontalPodAutoscaler should have
	MinReplicas int
}
