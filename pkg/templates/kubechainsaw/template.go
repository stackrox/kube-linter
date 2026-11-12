package kubechainsaw

import (
	"fmt"
	"strings"

	"github.com/ugiordan/kube-chainsaw/pkg/analyzer"
	kcModels "github.com/ugiordan/kube-chainsaw/pkg/models"
	"github.com/ugiordan/kube-chainsaw/pkg/suppression"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/k8sutil"
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
		Description: "Detects RBAC privilege escalation paths, dangerous verbs, sensitive resource access, and scope confusion using graph analysis. Rules that require NetworkPolicy, SecurityContextConstraints, Service, Route or Secret objects are not supported and are rejected at configuration time",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{
				objectkinds.ClusterRole,
				objectkinds.Role,
				objectkinds.ClusterRoleBinding,
				objectkinds.RoleBinding,
				objectkinds.ServiceAccount,
				objectkinds.Pod,
				objectkinds.Deployment,
				objectkinds.DaemonSet,
				objectkinds.StatefulSet,
				objectkinds.Job,
				objectkinds.CronJob,
				objectkinds.ReplicaSet,
			},
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

	// Suppressions are configuration, not per-context state: load them once
	// here so a missing or invalid file fails instantiation instead of
	// failing every lint run, and is not re-read per context.
	var sups []suppression.Suppression
	if p.SuppressionsFile != "" {
		var err error
		sups, err = suppression.LoadSuppressions(p.SuppressionsFile)
		if err != nil {
			return nil, fmt.Errorf("loading suppressions file %q: %w", p.SuppressionsFile, err)
		}
	}

	type contextAnalysis struct {
		findings    []kcModels.Finding
		initErr     []diagnostic.Diagnostic
		firstObject k8sutil.Object
	}

	// The check function is invoked once per object, and a run may lint
	// several contexts, so analysis results are cached per lint context.
	analyses := make(map[lintcontext.LintContext]*contextAnalysis)

	return func(lintCtx lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
		analysis, ok := analyses[lintCtx]
		if !ok {
			analysis = &contextAnalysis{firstObject: object.K8sObject}
			analyses[lintCtx] = analysis

			resources, err := convert.FromLintContext(lintCtx)
			if err != nil {
				analysis.initErr = []diagnostic.Diagnostic{{
					Message:  fmt.Sprintf("kube-chainsaw conversion error: %v", err),
					Severity: diagnostic.SeverityWarning,
				}}
			} else {
				findings := analyzer.Analyze(resources)
				findings = filterByRules(findings, p.Rules, p.ExcludeRules)
				findings = filterBySeverity(findings, p.MinSeverity)
				if len(sups) > 0 {
					findings = suppression.ApplySuppressions(findings, sups)
				}
				analysis.findings = findings
			}
		}

		// Conversion failures are reported once per run: the runner visits
		// each object once per run, so the error follows the object that
		// triggered the analysis and is reported again when a new run over
		// the same context revisits it.
		if analysis.initErr != nil {
			if object.K8sObject == analysis.firstObject {
				return analysis.initErr
			}
			return nil
		}

		return toDiagnostics(findingsForObject(analysis.findings, object))
	}, nil
}

func findingsForObject(findings []kcModels.Finding, object lintcontext.Object) []kcModels.Finding {
	info := object.K8sObject.GetObjectKind().GroupVersionKind()
	name := object.GetK8sObjectName().Name
	// Findings carry the normalized namespace, so an object with an omitted
	// namespace must be normalized the same way to match.
	namespace := convert.DefaultNamespace(info.Kind, object.GetK8sObjectName().Namespace)

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
			Severity: diagnostic.Severity(strings.ToLower(f.Severity.String())),
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

	minLevel, err := params.ParseSeverityLevel(minSeverity)
	if err != nil {
		// Should never happen: params are validated during template instantiation
		return findings
	}
	var filtered []kcModels.Finding
	for _, f := range findings {
		if int(f.Severity) >= minLevel {
			filtered = append(filtered, f)
		}
	}
	return filtered
}
