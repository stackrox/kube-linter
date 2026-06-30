package kubechainsaw

import (
	"fmt"
	"strings"
	"sync"

	"github.com/ugiordan/kube-chainsaw/pkg/analyzer"
	kcModels "github.com/ugiordan/kube-chainsaw/pkg/models"
	"github.com/ugiordan/kube-chainsaw/pkg/suppression"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/kubechainsaw/internal/convert"
	"golang.stackrox.io/kube-linter/pkg/templates/kubechainsaw/internal/params"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "RBAC Privilege Chain Analysis",
		Key:         "kube-chainsaw",
		Description: "Detects RBAC privilege escalation paths, dangerous verbs, sensitive resource access, and scope confusion using graph analysis",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.Any},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate:            params.WrapInstantiateFunc(analyze),
	})
}

func analyze(p params.Params) (check.Func, error) {
	// Perform custom validation on top of generated validation
	if err := p.ValidateCustom(); err != nil {
		return nil, err
	}

	var mu sync.Mutex
	var findings []kcModels.Finding
	var analyzed bool
	var initErr []diagnostic.Diagnostic

	return func(lintCtx lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
		mu.Lock()
		if !analyzed {
			resources, err := convert.FromLintContext(lintCtx)
			if err != nil {
				initErr = []diagnostic.Diagnostic{{
					Message:  fmt.Sprintf("kube-chainsaw conversion error: %v", err),
					Severity: "warning",
				}}
				analyzed = true
				mu.Unlock()
				return initErr
			}

			findings = analyzer.Analyze(resources)
			findings = filterByRules(findings, p.Rules, p.ExcludeRules)
			findings = filterBySeverity(findings, p.MinSeverity)

			if p.SuppressionsFile != "" {
				sups, supErr := suppression.LoadSuppressions(p.SuppressionsFile)
				if supErr != nil {
					initErr = []diagnostic.Diagnostic{{
						Message:  fmt.Sprintf("kube-chainsaw: failed to load suppressions: %v", supErr),
						Severity: "warning",
					}}
					analyzed = true
					mu.Unlock()
					return initErr
				}
				findings = suppression.ApplySuppressions(findings, sups)
			}
			analyzed = true
		}
		mu.Unlock()

		if initErr != nil {
			return initErr
		}

		kind := object.K8sObject.GetObjectKind().GroupVersionKind().Kind
		if !isRelevantKind(kind) {
			return nil
		}

		return toDiagnostics(findingsForObject(findings, object))
	}, nil
}

func isRelevantKind(kind string) bool {
	switch kind {
	case "ClusterRole", "Role", "ClusterRoleBinding", "RoleBinding",
		"ServiceAccount", "Pod", "Deployment", "DaemonSet",
		"StatefulSet", "Job", "CronJob", "ReplicaSet":
		return true
	}
	return false
}

func findingsForObject(findings []kcModels.Finding, object lintcontext.Object) []kcModels.Finding {
	info := object.K8sObject.GetObjectKind().GroupVersionKind()
	name := object.GetK8sObjectName().Name
	namespace := object.GetK8sObjectName().Namespace

	var matched []kcModels.Finding
	for _, f := range findings {
		if f.Suppressed {
			continue
		}
		if f.ResourceKind == info.Kind && f.ResourceName == name && f.ResourceNamespace == namespace {
			matched = append(matched, f)
		}
	}
	return matched
}

func toDiagnostics(findings []kcModels.Finding) []diagnostic.Diagnostic {
	if len(findings) == 0 {
		return nil
	}
	diags := make([]diagnostic.Diagnostic, len(findings))
	for i, f := range findings {
		diags[i] = diagnostic.Diagnostic{
			Message:  f.Title + ": " + f.Description,
			Severity: strings.ToLower(f.Severity.String()),
			Metadata: map[string]string{
				diagnostic.MetaKeyRuleID:      f.RuleID,
				diagnostic.MetaKeyFingerprint: f.Fingerprint,
				diagnostic.MetaKeyRemediation: f.Remediation,
			},
		}
	}
	return diags
}

func filterByRules(findings []kcModels.Finding, include, exclude []string) []kcModels.Finding {
	if len(include) == 0 && len(exclude) == 0 {
		return findings
	}

	includeSet := make(map[string]bool)
	for _, id := range include {
		includeSet[id] = true
	}
	excludeSet := make(map[string]bool)
	for _, id := range exclude {
		excludeSet[id] = true
	}

	var filtered []kcModels.Finding
	for _, f := range findings {
		if excludeSet[f.RuleID] {
			continue
		}
		if len(includeSet) > 0 && !includeSet[f.RuleID] {
			continue
		}
		filtered = append(filtered, f)
	}
	return filtered
}

func filterBySeverity(findings []kcModels.Finding, minSeverity string) []kcModels.Finding {
	if minSeverity == "" {
		return findings
	}

	minLevel := parseSeverityLevel(minSeverity)
	var filtered []kcModels.Finding
	for _, f := range findings {
		if int(f.Severity) >= minLevel {
			filtered = append(filtered, f)
		}
	}
	return filtered
}

func parseSeverityLevel(s string) int {
	switch strings.ToLower(s) {
	case "info", "note":
		return 0
	case "warning":
		return 1
	case "high", "error":
		return 2
	case "critical":
		return 3
	default:
		return 0
	}
}
