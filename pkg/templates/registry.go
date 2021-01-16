package templates

import (
	"fmt"
	"sort"

	"golang.stackrox.io/kube-linter/pkg/check"
)

var (
	allTemplates = make(map[string]check.Template)
)

// Register registers a template with the given name.
// Intended to be called at program init time.
func Register(t check.Template) {
	if _, ok := allTemplates[t.Key]; ok {
		panic(fmt.Sprintf("duplicate template: %v", t.Key))
	}
	allTemplates[t.Key] = t
}

// Get gets a template by name, returning a boolean indicating whether it was found.
func Get(name string) (check.Template, bool) {
	t, ok := allTemplates[name]
	return t, ok
}

// List returns all known templates, sorted by name.
func List() []check.Template {
	out := make([]check.Template, 0, len(allTemplates))
	for _, t := range allTemplates {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Key < out[j].Key
	})
	return out
}
