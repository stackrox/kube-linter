package params

// Params represents the params accepted by this template.
type Params struct {
	// A list of labels that will not cause the check to fail. For example, a label that is known to be populated at runtime by Kubernetes.
	IgnoredLabels []string `json:"ignoredLabels"`
}
