package config

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/internal/check"
)

// Config represents the config file format.
type Config struct {
	Checks []check.Check `json:"checks"`
}

// Load loads the config from the given path.
func Load(path string) (*Config, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "reading file")
	}
	var c Config
	if err := yaml.Unmarshal(contents, &c); err != nil {
		return nil, errors.Wrap(err, "unmarshaling config YAML")
	}
	return &c, nil
}