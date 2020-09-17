package configresolver

import (
	"golang.stackrox.io/kube-linter/internal/checkregistry"
	"golang.stackrox.io/kube-linter/internal/config"
	"golang.stackrox.io/kube-linter/internal/defaultchecks"
	"golang.stackrox.io/kube-linter/internal/errorhelpers"
	"golang.stackrox.io/kube-linter/internal/set"
)

// LoadCustomChecksInto loads the custom checks from the config into the check registry.
func LoadCustomChecksInto(cfg *config.Config, checkRegistry checkregistry.CheckRegistry) error {
	errorList := errorhelpers.NewErrorList("check registration")
	for i, check := range cfg.CustomChecks {
		if err := checkRegistry.Register(&cfg.CustomChecks[i]); err != nil {
			errorList.AddWrapf(err, "failed to register custom check %s", check.Name)
		}
	}
	return errorList.ToError()
}

// GetEnabledChecksAndValidate get the list of enabled checks based on the given config,
// and validates that they exist in the given checkRegistry.
func GetEnabledChecksAndValidate(cfg *config.Config, checkRegistry checkregistry.CheckRegistry) ([]string, error) {
	enabledChecks := set.NewStringSet()
	if !cfg.Checks.DoNotAutoAddDefaults {
		enabledChecks.AddAll(defaultchecks.List...)
	}
	for _, check := range cfg.CustomChecks {
		enabledChecks.Add(check.Name)
	}
	enabledChecks.AddAll(cfg.Checks.Include...)
	enabledChecks.RemoveAll(cfg.Checks.Exclude...)

	errorList := errorhelpers.NewErrorList("enabled checks validation")
	for check := range enabledChecks {
		if checkRegistry.Load(check) == nil {
			errorList.AddStringf("check %q not found", check)
		}
	}
	if err := errorList.ToError(); err != nil {
		return nil, err
	}
	return enabledChecks.AsSortedSlice(func(i, j string) bool {
		return i < j
	}), nil
}
