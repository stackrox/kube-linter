package params

// Params represents the params accepted by this template.
type Params struct {

	// The minimum number of replicas a deployment must have before anti-affinity is enforced on it
	MinReplicas int
}
