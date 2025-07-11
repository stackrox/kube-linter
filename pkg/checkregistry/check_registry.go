package checkregistry

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/instantiatedcheck"
)

// A CheckRegistry is a registry of checks.
// It is not thread-safe. It is anticipated that checks will all be registered ahead of time
// before calls to Load.
type CheckRegistry interface {
	Register(checks ...*config.Check) error
	Load(name string) *instantiatedcheck.InstantiatedCheck
}

type checkRegistry map[string]*instantiatedcheck.InstantiatedCheck

func (cr checkRegistry) Register(checks ...*config.Check) error {
	for _, c := range checks {
		instantiated, err := instantiatedcheck.ValidateAndInstantiate(c)
		if err != nil {
			return fmt.Errorf("invalid check %s: %w", c.Name, err)
		}
		if _, ok := cr[instantiated.Spec.Name]; ok {
			return fmt.Errorf("duplicate check name: %s", instantiated.Spec.Name)
		}
		cr[instantiated.Spec.Name] = instantiated
	}
	return nil
}

func (cr checkRegistry) Load(name string) *instantiatedcheck.InstantiatedCheck {
	return cr[name]
}

// New returns a ready-to-use, empty CheckRegistry.
func New() CheckRegistry {
	return make(checkRegistry)
}
