package params

// Params represents the params accepted by the lua script template.
type Params struct {
	// Script is the path to the Lua script file
	Script string

	// Inline is the inline Lua script content (alternative to Script)
	Inline string

	// Timeout is the execution timeout in seconds (default: 5)
	Timeout int
}
