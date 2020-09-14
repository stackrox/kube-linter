package templates

import (
	"fmt"

	"golang.stackrox.io/kube-linter/internal/check"
)

var (
	allTemplates = make(map[string]check.Template)
)

// Register registers a template with the given name.
// Intended to be called at program init time.
func Register(t check.Template) {
	if _, ok := allTemplates[t.Name]; ok {
		panic(fmt.Sprintf("duplicate template: %v", t.Name))
	}
	allTemplates[t.Name] = t
}

// Get gets a template by name, returning a boolean indicating whether it was found.
func Get(name string) (check.Template, bool) {
	t, ok := allTemplates[name]
	return t, ok
}
