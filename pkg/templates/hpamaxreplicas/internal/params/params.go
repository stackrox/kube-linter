package params

// Params represents the params accepted by this template.
type Params struct {

	// The maximum number of replicas a HorizontalPodAutoscaler should have
	MaxReplicas int
}
