package updateconfig

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"golang.stackrox.io/kube-linter/internal/errorhelpers"
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
)

const (
	templateKey = "update-configuration"
)

func parseIntOrString(data string) (*intstr.IntOrString, error) {
	val, err := strconv.Atoi(data)
	if err != nil {
		// This is not an integer.  Is it a valid string?
		if !strings.HasSuffix(data, "%") {
			return nil, fmt.Errorf("%s is not a valid string.  It does not end with %s", data, "%")
		}
		strLen := len(data)
		intVal := data[:strLen-1]
		val, err = strconv.Atoi(intVal)
		if err != nil {
			// Not going to try harder
			return nil, fmt.Errorf("unable to parse %s", data)
		}
		if val > 100 || val < 0 {
			return nil, fmt.Errorf("%s isn't a valid percent", data)
		}
	} else if val < 0 {
		return nil, fmt.Errorf("%d isn't a valid value", val)
	}
	converted := intstr.Parse(data)
	return &converted, nil

}

func compareIntOrString(maximal, minimal, actual *intstr.IntOrString) bool {
	if maximal == nil && minimal == nil {
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
	if maximal != nil {
		maxIsPercent := strings.Contains(maximal.String(), "%")
		if actualIsPercent != maxIsPercent {
			return false
		}
		maxVal, err := intstr.GetValueFromIntOrPercent(maximal, 100, false)
		if err != nil {
			return false
		}
		if actualVal > maxVal {
			return false
		}
	}
	if minimal != nil {
		minIsPercent := strings.Contains(minimal.String(), "%")
		if actualIsPercent != minIsPercent {
			return false
		}
		minVal, err := intstr.GetValueFromIntOrPercent(minimal, 100, false)
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
			var maxPodsUnavailable, minPodsUnavailable, maxSurge, minSurge *intstr.IntOrString
			errorList := errorhelpers.NewErrorList("configuration verification")
			compiledRegex, err := regexp.Compile(p.StrategyTypeRegex)
			if err != nil {
				errorList.AddWrapf(err, "invalid regex %s", p.StrategyTypeRegex)
			}
			if len(p.MaxPodsUnavailable) > 0 {
				maxPodsUnavailable, err = parseIntOrString(p.MaxPodsUnavailable)
				if err != nil {
					errorList.AddWrapf(err, "invalid MaxPodsUnavailable %s", p.MaxPodsUnavailable)
				}
			}
			if len(p.MinPodsUnavailable) > 0 {
				minPodsUnavailable, err = parseIntOrString(p.MinPodsUnavailable)
				if err != nil {
					errorList.AddWrapf(err, "invalid MinPodsUnavailable %s", p.MinPodsUnavailable)
				}
			}
			if len(p.MaxSurge) > 0 {
				maxSurge, err = parseIntOrString(p.MaxSurge)
				if err != nil {
					errorList.AddWrapf(err, "invalid MaxSurge %s", p.MaxSurge)
				}
			}
			if len(p.MinSurge) > 0 {
				minSurge, err = parseIntOrString(p.MinSurge)
				if err != nil {
					errorList.AddWrapf(err, "invalid MinSurge %s", p.MinSurge)
				}
			}
			if err := errorList.ToError(); err != nil {
				return nil, err
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
					if !compareIntOrString(maxPodsUnavailable, minPodsUnavailable, strategy.MaxUnavailable) {
						minStr := fmt.Sprintf("at least %s", p.MinPodsUnavailable)
						maxStr := fmt.Sprintf("no more than %s", p.MaxPodsUnavailable)
						msg := fmt.Sprintf("object has a max unavailable of %s but %s is required", strategy.MaxUnavailable.String(),
							conditional(len(p.MinPodsUnavailable) > 0, minStr, len(p.MaxPodsUnavailable) > 0, maxStr, " and "))
						newD := diagnostic.Diagnostic{Message: msg}
						diagnostics = append(diagnostics, newD)
					}
				}
				if strategy.MaxSurgeExists {
					if !compareIntOrString(maxSurge, minSurge, strategy.MaxSurge) {
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
