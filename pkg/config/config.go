package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// ChecksConfig is the config that determines which checks to run.
type ChecksConfig struct {
	// AddAllBuiltIn, if set, adds all built-in checks. This allows users to
	// explicitly opt-out of checks that are not relevant using Exclude.
	// +flagName=add-all-built-in
	AddAllBuiltIn bool `json:"addAllBuiltIn"`
	// DoNotAutoAddDefaults, if set, prevents the automatic addition of default checks.
	// +flagName=do-not-auto-add-defaults
	DoNotAutoAddDefaults bool `json:"doNotAutoAddDefaults"`
	// Exclude is a list of check names to exclude.
	// +flagName=exclude
	Exclude []string `json:"exclude"`
	// Include is a list of check names to include. If a check is in both Include and Exclude,
	// Exclude wins.
	// +flagName=include
	Include []string `json:"include"`
	// IgnorePaths is a list of path to ignore from applying checks
	// +flagName=ignore-paths
	IgnorePaths []string `json:"ignorePaths"`
}

// Config represents the config file format.
type Config struct {
	// +flagName=-
	CustomChecks []Check      `json:"customChecks,omitempty"`
	Checks       ChecksConfig `json:"checks,omitempty"`
}

// Defines the list of default config filenames to check if parameter isn't passed in
var defaultConfigFilenames = [...]string{".kube-linter.yaml", ".kube-linter.yml"}

// Get info on config file if it exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// Load loads the config from the given path.
func Load(v *viper.Viper, configPath string) (Config, error) {
	if configPath == "" {
		for _, p := range defaultConfigFilenames {
			if fileExists(p) {
				configPath = p
				break
			}
		}
	}

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
	err := v.Unmarshal(&conf, viper.DecoderConfigOption(func(config *mapstructure.DecoderConfig) {
		config.TagName = "json"
	}))
	if err != nil {
		return Config{}, errors.Wrap(err, "unmarshalling config File")
	}
	return conf, nil
}
