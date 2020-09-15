package run

import (
	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/internal/config"
	"golang.stackrox.io/kube-linter/internal/diagnostic"
	"golang.stackrox.io/kube-linter/internal/lintcontext"
)

// Result represents the result from a run of the linter.
type Result struct {
	Reports []diagnostic.WithContext
}

// Run runs the linter on the given context, with the given config.
func Run(lintCtx *lintcontext.LintContext, cfg *config.Config) (Result, error) {
	var instantiatedChecks []*instantiatedCheck
	for i, check := range cfg.Checks {
		instantiatedCheck, err := validateAndInstantiate(&cfg.Checks[i])
		if err != nil {
			return Result{}, errors.Wrapf(err, "invalid check %q", check.Name)
		}
		instantiatedChecks = append(instantiatedChecks, instantiatedCheck)
	}

	var result Result
	for _, obj := range lintCtx.Objects {
		for _, check := range instantiatedChecks {
			if obj.K8sObject == nil {
				continue
			}
			if !check.Matcher.Matches(obj.K8sObject.GetObjectKind().GroupVersionKind()) {
				continue
			}
			diagnostics := check.Func(lintCtx, obj)
			for _, d := range diagnostics {
				result.Reports = append(result.Reports, diagnostic.WithContext{
					Diagnostic: d,
					Check:      check.Name,
					Object:     obj,
				})
			}
		}
	}
	return result, nil
}
