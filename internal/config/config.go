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
	AddAllBuiltIn bool `json:"addAllBuiltIn"`
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
func Load(v *viper.Viper) (Config, error) {
	configFile := v.GetString("config")

	if configFile != "" {
		filename := filepath.Base(configFile)
		ext := filepath.Ext(configFile)
		configPath := filepath.Dir(configFile)

		v.SetConfigType(strings.TrimPrefix(ext, "."))
		v.SetConfigName(strings.TrimSuffix(filename, ext))
		v.AddConfigPath(configPath)
		if err := v.ReadInConfig(); err != nil {
			return Config{}, errors.Wrap(err, "reading file")
		}
	}

	var conf Config
	err := v.Unmarshal(&conf)
	if err != nil {
		return Config{}, errors.Wrap(err, "unmarshalling config YAML")
	}
	return conf, nil
}
