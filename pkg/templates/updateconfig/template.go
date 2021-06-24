package updateconfig

import (
	"fmt"
	"regexp"
	"strings"

	"golang.stackrox.io/kube-linter/internal/stringutils"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/updateconfig/internal/params"

	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/pkg/errors"
)

const (
	templateKey = "update-configuration"
)

func compareIntOrString(max, min string, actual *intstr.IntOrString) bool {
	if len(max) == 0 && len(min) == 0 {
		return true
	}
	if actual == nil {
		return false
	}
	actualVal, err := intstr.GetValueFromIntOrPercent(actual, 100, false)
	if err != nil {
		return false
	}
	actualIsPercent := strings.Contains(actual.String(), "%")
	if len(max) > 0 {
		maxIntOrString := intstr.Parse(max)
		maxIsPercent := strings.Contains(maxIntOrString.String(), "%")
		if actualIsPercent != maxIsPercent {
			return false
		}
		maxVal, err := intstr.GetValueFromIntOrPercent(&maxIntOrString, 100, false)
		if err != nil {
			return false
		}
		if actualVal > maxVal {
			return false
		}
	}
	if len(min) > 0 {
		minIntOrString := intstr.Parse(min)
		minIsPercent := strings.Contains(minIntOrString.String(), "%")
		if actualIsPercent != minIsPercent {
			return false
		}
		minVal, err := intstr.GetValueFromIntOrPercent(&minIntOrString, 100, false)
		if err != nil {
			return false
		}
		if actualVal < minVal {
			return false
		}
	}
	return true
}

func conditional(firstCond bool, firstStr string, secondCond bool, secondStr, bothStr string) string {
	msg := ""
	if firstCond {
		msg = firstStr
	}
	if firstCond && secondCond {
		msg += bothStr
	}
	if secondCond {
		msg += secondStr
	}
	return msg

}

func needsRollingUpdateDefinition(p params.Params) bool {
	isRolling, _ := regexp.MatchString("Rolling", p.StrategyTypeRegex)
	return isRolling && (len(p.MinPodsUnavailable) > 0 || len(p.MaxPodsUnavailable) > 0 ||
		len(p.MinSurge) > 0 || len(p.MaxSurge) > 0)
}

func init() {
	templates.Register(check.Template{
		HumanName:   "Update configuration",
		Key:         templateKey,
		Description: "Flag configurations that do not meet the specified update configuration",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			compiledRegex, err := regexp.Compile(p.StrategyTypeRegex)
			if err != nil {
				return nil, errors.Wrapf(err, "invalid regex %s", p.StrategyTypeRegex)
			}
			return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				var diagnostics []diagnostic.Diagnostic

				strategy, found := extract.UpdateStrategy(object.K8sObject)
				if !found {
					return nil
				}
				if !strategy.TypeExists {
					return nil
				}
				if !compiledRegex.MatchString(strategy.Type) {
					newD := diagnostic.Diagnostic{
						Message: fmt.Sprintf("object has %s strategy type but must match regex %s",
							stringutils.Ternary(len(strategy.Type) > 0, strategy.Type, "no"), p.StrategyTypeRegex)}
					diagnostics = append(diagnostics, newD)
				}
				if !strategy.RollingConfigExists {
					return nil
				}
				if needsRollingUpdateDefinition(p) && !strategy.RollingConfigValid {
					newD := diagnostic.Diagnostic{Message: "object has no rolling update parameters defined"}
					diagnostics = append(diagnostics, newD)
				}
				if strategy.MaxUnavailableExists {
					if !compareIntOrString(p.MaxPodsUnavailable, p.MinPodsUnavailable, strategy.MaxUnavailable) {
						minStr := fmt.Sprintf("at least %s", p.MinPodsUnavailable)
						maxStr := fmt.Sprintf("no more than %s", p.MaxPodsUnavailable)
						msg := fmt.Sprintf("object has a max unavailable of %s but %s is required", strategy.MaxUnavailable.String(),
							conditional(len(p.MinPodsUnavailable) > 0, minStr, len(p.MaxPodsUnavailable) > 0, maxStr, " and "))
						newD := diagnostic.Diagnostic{Message: msg}
						diagnostics = append(diagnostics, newD)
					}
				}
				if strategy.MaxSurgeExists {
					if !compareIntOrString(p.MaxSurge, p.MinSurge, strategy.MaxSurge) {
						minStr := fmt.Sprintf("at least %s", p.MinSurge)
						maxStr := fmt.Sprintf("no more than %s", p.MaxSurge)
						msg := fmt.Sprintf("object has a max surge of %s but %s is required", strategy.MaxSurge.String(),
							conditional(len(p.MinSurge) > 0, minStr, len(p.MaxSurge) > 0, maxStr, " and "))
						newD := diagnostic.Diagnostic{Message: msg}
						diagnostics = append(diagnostics, newD)
					}
				}
				if len(diagnostics) == 0 {
					return nil
				}
				return diagnostics
			}, nil
		}),
	})
}
