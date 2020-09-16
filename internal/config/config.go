package config

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/internal/check"
)

// ChecksConfig is the config that determines which checks to run.
type ChecksConfig struct {
	// DoNotAutoAddDefaults, if set, prevents the automatic addition of default checks.
	DoNotAutoAddDefaults bool `json:"doNotAutoAddDefaults"`
	// Exclude is a list of check names to exclude.
	Exclude []string `json:"exclude"`
	// Include is a list of check names to include. If a check is in both Include and Exclude,
	// Exclude wins.
	Include []string `json:"include"`
}

// Config represents the config file format.
type Config struct {
	CustomChecks []check.Check `json:"customChecks,omitempty"`
	Checks       ChecksConfig  `json:"checks,omitempty"`
}

// Load loads the config from the given path.
func Load(path string) (Config, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, errors.Wrap(err, "reading file")
	}
	var c Config
	if err := yaml.Unmarshal(contents, &c); err != nil {
		return Config{}, errors.Wrap(err, "unmarshaling config YAML")
	}
	return c, nil
}
