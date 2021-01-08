package run

import (
	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/internal/checkregistry"
	"golang.stackrox.io/kube-linter/internal/diagnostic"
	"golang.stackrox.io/kube-linter/internal/ignore"
	"golang.stackrox.io/kube-linter/internal/instantiatedcheck"
	"golang.stackrox.io/kube-linter/internal/lintcontext"
)

// Result represents the result from a run of the linter.
type Result struct {
	Reports []diagnostic.WithContext
}

// ReportObject describes underlying k8s object
type ReportObject struct {
	Namespace string
	Name      string
	Version   string
}

// ReportLine describes individual line in linting report
type ReportLine struct {
	Path        string
	Object      ReportObject
	Message     string
	Check       string
	Remediation string
}

// GenerateReport consolidates the report data.
func (r *Result) GenerateReport() []ReportLine {
	result := []ReportLine{}

	for _, report := range r.Reports {

		line := ReportLine{
			report.Object.Metadata.FilePath,
			ReportObject{
				report.Object.K8sObject.GetNamespace(),
				report.Object.K8sObject.GetName(),
				report.Object.K8sObject.GetObjectKind().GroupVersionKind().String()},
			report.Diagnostic.Message,
			report.Check,
			report.Remediation}

		result = append(result, line)
	}

	return result
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
	return result, nil
}
