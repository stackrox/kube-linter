package builtinchecks

import (
	"github.com/ghodss/yaml"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/checkregistry"
)

var (
	box = packr.NewBox("./yamls")
)

// LoadInto loads built-in checks into the registry.
func LoadInto(registry checkregistry.CheckRegistry) error {
	checks, err := List()
	if err != nil {
		return err
	}
	for _, chk := range checks {
		if err := registry.Register(&chk); err != nil {
			return errors.Wrapf(err, "registering default check %s", chk.Name)
		}
	}
	return nil
}

// List lists built-in checks.
func List() ([]check.Check, error) {
	var out []check.Check
	for _, fileName := range box.List() {
		contents, err := box.Find(fileName)
		if err != nil {
			return nil, errors.Wrapf(err, "loading default check from %s", fileName)
		}
		var chk check.Check
		if err := yaml.Unmarshal(contents, &chk); err != nil {
			return nil, errors.Wrapf(err, "unmarshaling default check from %s", fileName)
		}
		out = append(out, chk)
	}
	return out, nil
}