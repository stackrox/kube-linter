package params

// Params represents the params accepted by this template.
type Params struct {
	// Recursive determines whether to check keys recursively at all nesting levels.
	// Default is true.
	Recursive bool

	// IgnoredKeys is a list of keys whose values are not checked for sorting, at any nesting level,
	// for example "labels" and "annotations". The order of the ignored key itself among its siblings is still checked.
	// +noregex
	// +notnegatable
	IgnoredKeys []string `json:"ignoredKeys"`
}
