package run

import (
	"time"

	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/internal/version"
	"golang.stackrox.io/kube-linter/pkg/checkregistry"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/ignore"
	"golang.stackrox.io/kube-linter/pkg/instantiatedcheck"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
)

type checkStatus string

const (
	checksPassed checkStatus = "Passed"
	checksFailed checkStatus = "Failed"
)

// Result represents the result from a run of the linter.
type Result struct {
	Reports []diagnostic.WithContext
	Summary Summary
}

// Summary holds information about the linter run overall.
type Summary struct {
	ChecksStatus      checkStatus
	CheckTime         time.Time
	KubeLinterVersion string
}

// Run runs the linter on the given context, with the given config.
func Run(lintCtxs []lintcontext.LintContext, registry checkregistry.CheckRegistry, checks []string) (Result, error) {
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
				if ignore.ObjectForCheck(obj.K8sObject.GetAnnotations(), check.Spec.Name) {
					continue
				}
				diagnostics := check.Func(lintCtx, obj)
				for _, d := range diagnostics {
					result.Reports = append(result.Reports, diagnostic.WithContext{
						Diagnostic:  d,
						Check:       check.Spec.Name,
						Remediation: check.Spec.Remediation,
						Object:      obj,
					})
				}
			}
		}
	}

	if len(result.Reports) > 0 {
		result.Summary.ChecksStatus = checksFailed
	} else {
		result.Summary.ChecksStatus = checksPassed
	}
	result.Summary.CheckTime = time.Now().UTC()
	result.Summary.KubeLinterVersion = version.Get()

	return result, nil
}
