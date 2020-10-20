package check

// A Check represents a single check. It is serializable.
type Check struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Scope       *ObjectKindsDesc       `json:"scope"`
	Template    string                 `json:"template"`
	Params      map[string]interface{} `json:"params,omitempty"`
}
