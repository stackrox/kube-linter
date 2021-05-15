package params

// Params represents the params accepted by this template.
type Params struct {
	// An array of regular expressions specifying system directories to be mounted on containers. e.g. ^/usr$ for /usr
	// +notnegatable
	Dirs []string `json:"dirs"`
}
