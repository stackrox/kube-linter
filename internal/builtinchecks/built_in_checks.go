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
	for _, fileName := range box.List() {
		contents, err := box.Find(fileName)
		if err != nil {
			return errors.Wrapf(err, "loading default check from %s", fileName)
		}
		var chk check.Check
		if err := yaml.Unmarshal(contents, &chk); err != nil {
			return errors.Wrapf(err, "unmarshaling default check from %s", fileName)
		}
		if err := registry.Register(&chk); err != nil {
			return errors.Wrapf(err, "registering default check from %s", fileName)
		}
	}
	return nil
}