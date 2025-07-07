package builtinchecks

import (
	"embed"
	"fmt"
	"path/filepath"
	"sync"

	"golang.stackrox.io/kube-linter/pkg/checkregistry"
	"golang.stackrox.io/kube-linter/pkg/config"
	"sigs.k8s.io/yaml"
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
	for i := range checks {
		if err := registry.Register(&checks[i]); err != nil {
			return fmt.Errorf("registering default check %s: %w", checks[i].Name, err)
		}
	}
	return nil
}

// List lists built-in checks.
func List() ([]config.Check, error) {
	loadOnce.Do(func() {
		fileEntries, err := yamlFiles.ReadDir("yamls")
		if err != nil {
			loadErr = fmt.Errorf("reading embedded yaml files: %w", err)
			return
		}
		for _, entry := range fileEntries {
			if entry.IsDir() || filepath.Ext(entry.Name()) != ".yaml" {
				loadErr = fmt.Errorf("found unexpected entry %s in yamls directory", entry.Name())
				return
			}
			// Do NOT use filepath.Join here, because embed always uses `/` as the separator,
			// irrespective of the OS we're running.
			contents, err := yamlFiles.ReadFile(fmt.Sprintf("yamls/%s", entry.Name()))
			if err != nil {
				loadErr = fmt.Errorf("loading file %s: %w", entry.Name(), err)
				return
			}
			var chk config.Check
			if err := yaml.Unmarshal(contents, &chk); err != nil {
				loadErr = fmt.Errorf("unmarshalling default check from %s: %w", entry.Name(), err)
				return
			}
			list = append(list, chk)
		}
	})
	if loadErr != nil {
		return nil, fmt.Errorf("UNEXPECTED: failed to load built-in checks: %w", loadErr)
	}
	return list, nil
}
