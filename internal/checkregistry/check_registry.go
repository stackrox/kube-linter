package checkregistry

import (
	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/instantiatedcheck"
)

// A CheckRegistry is a registry of checks.
// It is not thread-safe. It is anticipated that checks will all be registered ahead of time
// before calls to Load.
type CheckRegistry interface {
	Register(checks ...*check.Check) error
	Load(name string) *instantiatedcheck.InstantiatedCheck
}

type checkRegistry map[string]*instantiatedcheck.InstantiatedCheck

func (cr checkRegistry) Register(checks ...*check.Check) error {
	for _, c := range checks {
		instantiated, err := instantiatedcheck.ValidateAndInstantiate(c)
		if err != nil {
			return errors.Wrapf(err, "invalid check %s", c.Name)
		}
		if _, ok := cr[instantiated.Name]; ok {
			return errors.Errorf("duplicate check name: %s", instantiated.Name)
		}
		cr[instantiated.Name] = instantiated
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
