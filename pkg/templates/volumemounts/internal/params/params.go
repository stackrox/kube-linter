package params

// Params represents the params accepted by this template.
type Params struct {
	// An array of names of sensitive system directories to be mounted on containers
	// +noregex
	// +notnegatable
	SensitiveSysDirs []string `json:"sensitiveSysDirs"`
}
