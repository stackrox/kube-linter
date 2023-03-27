package configresolver

import (
	"golang.stackrox.io/kube-linter/internal/defaultchecks"
	"golang.stackrox.io/kube-linter/internal/errorhelpers"
	"golang.stackrox.io/kube-linter/internal/set"
	"golang.stackrox.io/kube-linter/pkg/builtinchecks"
	"golang.stackrox.io/kube-linter/pkg/checkregistry"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/pathutil"
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
		enabledChecks.AddAll(defaultchecks.List.AsSlice()...)
	}
	if cfg.Checks.AddAllBuiltIn {
		builtInChecks, err := builtinchecks.List()
		if err != nil {
			return nil, err
		}
		for _, check := range builtInChecks {
			enabledChecks.Add(check.Name)
		}
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

// GetIgnorePaths loads the paths from the config into the check registry.
func GetIgnorePaths(cfg *config.Config) ([]string, error) {
	errorList := errorhelpers.NewErrorList("check ignore paths")
	ignorePaths := set.NewStringSet()
	for _, path := range cfg.Checks.IgnorePaths {
		res, err := pathutil.GetAbsolutPath(path)
		if err != nil {
			errorList.AddError(err)
			continue
		}
		ignorePaths.AddAll(res)
	}

	if err := errorList.ToError(); err != nil {
		return nil, err
	}
	return ignorePaths.AsSortedSlice(func(i, j string) bool {
		return i < j
	}), nil
}
