package builtinchecks

import (
	"embed"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/pkg/checkregistry"
	"golang.stackrox.io/kube-linter/pkg/config"
)

var (
	//go:embed yamls
	yamlFiles embed.FS

	loadOnce sync.Once
	list     []config.Check
	loadErr  error
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
func List() ([]config.Check, error) {
	loadOnce.Do(func() {
		fileEntries, err := yamlFiles.ReadDir("yamls")
		if err != nil {
			loadErr = errors.Wrap(err, "reading embedded yaml files")
			return
		}
		for _, entry := range fileEntries {
			if entry.IsDir() || filepath.Ext(entry.Name()) != ".yaml" {
				loadErr = errors.Errorf("found unexpected entry %s in yamls directory", entry.Name())
				return
			}
			// Do NOT use filepath.Join here, because embed always uses `/` as the separator,
			// irrespective of the OS we're running.
			contents, err := yamlFiles.ReadFile(fmt.Sprintf("yamls/%s", entry.Name()))
			if err != nil {
				loadErr = errors.Wrapf(err, "loading file %s", entry.Name())
				return
			}
			var chk config.Check
			if err := yaml.Unmarshal(contents, &chk); err != nil {
				loadErr = errors.Wrapf(err, "unmarshalling default check from %s", entry.Name())
				return
			}
			list = append(list, chk)
		}
	})
	if loadErr != nil {
		return nil, errors.Wrap(loadErr, "UNEXPECTED: failed to load built-in checks")
	}
	return list, nil
}
