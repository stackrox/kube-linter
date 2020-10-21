package run

import (
	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/internal/checkregistry"
	"golang.stackrox.io/kube-linter/internal/diagnostic"
	"golang.stackrox.io/kube-linter/internal/instantiatedcheck"
	"golang.stackrox.io/kube-linter/internal/lintcontext"
)

// Result represents the result from a run of the linter.
type Result struct {
	Reports []diagnostic.WithContext
}

// Run runs the linter on the given context, with the given config.
func Run(lintCtxs []*lintcontext.LintContext, registry checkregistry.CheckRegistry, checks []string) (Result, error) {

	instantiatedChecks := make([]*instantiatedcheck.InstantiatedCheck, 0, len(checks))
	for _, checkName := range checks {
		instantiatedCheck := registry.Load(checkName)
		if instantiatedCheck == nil {
			return Result{}, errors.Errorf("check %q not found", checkName)
		}
		instantiatedChecks = append(instantiatedChecks, instantiatedCheck)
	}

	var result Result
	for _, lintCtx := range lintCtxs {
		for _, obj := range lintCtx.Objects() {
			for _, check := range instantiatedChecks {
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
	}
	return result, nil
}
