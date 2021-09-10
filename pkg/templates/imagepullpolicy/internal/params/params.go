package params

// Params represents the params accepted by this template.
type Params struct {
	// list of forbidden image pull policy
	// +noregex
	// +notnegatable
	// +enum=Always
	// +enum=IfNotPresent
	// +enum=Never
	ForbiddenPolicies []string
}
