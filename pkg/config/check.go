package config

// A Check represents a single check. It is serializable.
type Check struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Remediation string                 `json:"remediation"`
	Scope       *ObjectKindsDesc       `json:"scope"`
	Template    string                 `json:"template"`
	Params      map[string]interface{} `json:"params,omitempty"`
}

// ObjectKindsDesc describes a list of supported object kinds for a check template.
type ObjectKindsDesc struct {
	ObjectKinds []string `json:"objectKinds"`
}
