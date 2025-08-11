package run

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/open-policy-agent/opa/v1/ast"
	"github.com/open-policy-agent/opa/v1/rego"
	"github.com/open-policy-agent/opa/v1/util"
	"golang.stackrox.io/kube-linter/internal/version"
	"golang.stackrox.io/kube-linter/pkg/checkregistry"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/ignore"
	"golang.stackrox.io/kube-linter/pkg/instantiatedcheck"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
)

// CheckStatus is enum type.
type CheckStatus string

const (
	// ChecksPassed means no lint errors found.
	ChecksPassed CheckStatus = "Passed"
	// ChecksFailed means lint errors were found.
	ChecksFailed CheckStatus = "Failed"
)

// Result represents the result from a run of the linter.
type Result struct {
	Checks  []config.Check
	Reports []diagnostic.WithContext
	Summary Summary
}

// Summary holds information about the linter run overall.
type Summary struct {
	ChecksStatus      CheckStatus
	CheckEndTime      time.Time
	KubeLinterVersion string
}

type input struct {
	Object interface{}            `json:"object"`
	Params map[string]interface{} `json:"params"`
}

// Run runs the linter on the given context, with the given config.
func Run(lintCtxs []lintcontext.LintContext, registry checkregistry.CheckRegistry, checks []string, repo string) (Result, error) {
	var result Result

	instantiatedChecks := make([]*instantiatedcheck.InstantiatedCheck, 0, len(checks))
	for _, checkName := range checks {
		instantiatedCheck := registry.Load(checkName)
		if instantiatedCheck == nil {
			return Result{}, fmt.Errorf("check %q not found", checkName)
		}
		instantiatedChecks = append(instantiatedChecks, instantiatedCheck)
		result.Checks = append(result.Checks, instantiatedCheck.Spec)
	}

	modules := []func(*rego.Rego){
		rego.SetRegoVersion(ast.RegoV1),
		rego.Load([]string{repo}, nil),
	}

	params := map[string]interface{}{}

	for _, lintCtx := range lintCtxs {
		for _, obj := range lintCtx.Objects() {
			for _, check := range instantiatedChecks {
				if !check.Matcher.Matches(obj.K8sObject.GetObjectKind().GroupVersionKind()) {
					continue
				}
				if ignore.ObjectForCheck(obj.K8sObject.GetAnnotations(), check.Spec.Name) {
					continue
				}

				template := strings.ReplaceAll(check.Spec.Template, "-", "")
				if template != "latesttag" {
					continue
				}

				params[template] = check.Spec.Params

				in := input{
					Object: obj.K8sObject,
					Params: params,
				}

				inMap := util.MustUnmarshalJSON(util.MustMarshalJSON(in))

				eval := rego.New(
					append([]func(*rego.Rego){
						rego.Query(fmt.Sprintf("data.kubelinter.template.%s.deny", template)),
						rego.Input(inMap),
					},
						modules...)...,
				)

				rs, err := eval.Eval(context.Background())

				if err != nil {
					result.Reports = append(result.Reports, diagnostic.WithContext{
						Diagnostic: diagnostic.Diagnostic{
							Message: err.Error(),
						},
						Check:       check.Spec.Name,
						Remediation: check.Spec.Remediation,
						Object:      obj,
					})
				}

				fromResult, err := messagesFromResult(rs)
				if err != nil {
					result.Reports = append(result.Reports, diagnostic.WithContext{
						Diagnostic: diagnostic.Diagnostic{
							Message: err.Error(),
						},
						Check:       check.Spec.Name,
						Remediation: check.Spec.Remediation,
						Object:      obj,
					})
				}
				for _, d := range fromResult {
					result.Reports = append(result.Reports, diagnostic.WithContext{
						Diagnostic:  diagnostic.Diagnostic{Message: d},
						Check:       check.Spec.Name,
						Remediation: check.Spec.Remediation,
						Object:      obj,
					})
				}
			}
		}
	}

	if len(result.Reports) > 0 {
		result.Summary.ChecksStatus = ChecksFailed
	} else {
		result.Summary.ChecksStatus = ChecksPassed
	}
	result.Summary.CheckEndTime = time.Now().UTC()
	result.Summary.KubeLinterVersion = version.Get()

	return result, nil
}

func messagesFromResult(rs rego.ResultSet) ([]string, error) {
	var messages []string
	for _, result := range rs {
		for _, r := range result.Expressions {
			msgs, ok := r.Value.([]interface{})
			if !ok {
				return nil, fmt.Errorf("unexpected value %v", r.Value)
			}
			for _, v := range msgs {
				str, ok := v.(string)
				if !ok {
					return nil, fmt.Errorf("unexpected value %v", v)
				}
				messages = append(messages, str)
			}
		}
	}
	return messages, nil
}
