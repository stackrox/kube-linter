package params

// Params represents the params accepted by this template.
type Params struct {
	// An array of unsafe system controls
	// +noregex
	// +notnegatable
	UnsafeSysCtls []string `json:"unsafeSysCtls"`
}
