package check

// A Check represents a single check. It is serializable.
type Check struct {
	Name     string            `json:"name"`
	Scope    *ObjectKindsDesc  `json:"scope"`
	Template string            `json:"template"`
	Params   map[string]string `json:"params,omitempty"`
}
