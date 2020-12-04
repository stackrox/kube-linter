package params

// Params represents the params accepted by this template.
type Params struct {

	// Comma (ex: "NET_ADMIN,SYS_TIME") separated list of capabilities that are required to
	// be dropped by containers.
	RequiredDrops string `json:"requiredDrops"`

	// Comma (ex: "NET_ADMIN,SYS_TIME") separated list of capabilities that are forbidden to
	// be added to containers.
	ForbiddenAdds string `json:"forbiddenAdds"`
}
