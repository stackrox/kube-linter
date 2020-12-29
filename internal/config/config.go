package config

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.stackrox.io/kube-linter/internal/check"
)

// ChecksConfig is the config that determines which checks to run.
type ChecksConfig struct {
	// AddAllBuiltIn, if set, adds all built-in checks. This allows users to
	// explicitly opt-out of checks that are not relevant using Exclude.
	AddAllBuiltIn bool `json:"addAllBuiltIn" mapstructure:"add-all-built-in"`
	// DoNotAutoAddDefaults, if set, prevents the automatic addition of default checks.
	DoNotAutoAddDefaults bool `json:"doNotAutoAddDefaults" mapstructure:"do-not-auto-add-defaults"`
	// Exclude is a list of check names to exclude.
	Exclude []string `json:"exclude" mapstructure:"exclude"`
	// Include is a list of check names to include. If a check is in both Include and Exclude,
	// Exclude wins.
	Include []string `json:"include" mapstructure:"include"`
}

// Config represents the config file format.
type Config struct {
	// +viper=exclude
	CustomChecks []check.Check `json:"customChecks,omitempty" mapstructure:"customChecks,omitempty"`
	Checks       ChecksConfig  `json:"checks,omitempty" mapstructure:"checks"`
}

// Load loads the config from the given path.
func Load(v *viper.Viper, configPath string) (Config, error) {

	if configPath != "" {
		filename := filepath.Base(configPath)
		ext := filepath.Ext(configPath)
		path := filepath.Dir(configPath)

		v.SetConfigName(strings.TrimSuffix(filename, ext))
		v.AddConfigPath(path)
		if err := v.ReadInConfig(); err != nil {
			return Config{}, errors.Wrap(err, "reading file")
		}
	}

	var conf Config
	err := v.Unmarshal(&conf)
	if err != nil {
		return Config{}, errors.Wrap(err, "unmarshalling config File")
	}
	return conf, nil
}
