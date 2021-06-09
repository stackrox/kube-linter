package params

// Params represents the params accepted by this template.
type Params struct {
	// Set to true to flag the roles that are referenced in bindings but not found in the context
	FlagRolesNotFound bool `json:"flagRolesNotFound"`
	// An array of regular expressions specifying resources. e.g. ^secrets$ for secrets and ^*$ for any resources
	// +notnegatable
	Resources []string `json:"resources"`
	// An array of regular expressions specifying verbs. e.g. ^create$ for create and ^*$ for any k8s verbs
	// +notnegatable
	Verbs []string `json:"verbs"`
}
