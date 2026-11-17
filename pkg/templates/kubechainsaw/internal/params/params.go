package params

import (
	"fmt"
	"strings"

	"github.com/ugiordan/kube-chainsaw/pkg/analyzer"
)

// Params holds the configuration for the kube-chainsaw template.
type Params struct {
	Rules            []string `json:"rules"`
	ExcludeRules     []string `json:"excludeRules"`
	MinSeverity      string   `json:"minSeverity"`
	SuppressionsFile string   `json:"suppressionsFile"`
}

// unsupportedRuleIDs are rules that analyze resource kinds the kube-linter
// integration does not load (NetworkPolicy, SecurityContextConstraints,
// Service, Ingress, Route and Secret objects). They are rejected at
// instantiation because selecting them could never produce findings here.
var unsupportedRuleIDs = map[string]bool{
	analyzer.RuleNetworkPolicyIngress:      true,
	analyzer.RuleNetworkPolicyEgress:       true,
	analyzer.RulePermissiveSCC:             true,
	analyzer.RuleExternalService:           true,
	analyzer.RuleUnencryptedRoute:          true,
	analyzer.RuleBroadExternalRoute:        true,
	analyzer.RuleCredentialSecret:          true,
	analyzer.RuleServiceAccountTokenSecret: true,
}

// SupportedRuleIDs returns the rule IDs this integration can evaluate.
func SupportedRuleIDs() []string {
	all := analyzer.KnownRuleIDs()
	supported := make([]string, 0, len(all))
	for _, id := range all {
		if !unsupportedRuleIDs[id] {
			supported = append(supported, id)
		}
	}
	return supported
}

func validateRuleID(id string, known map[string]bool, paramName string) error {
	if !known[id] {
		return fmt.Errorf("unknown %s ID %q; valid IDs: %s", paramName, id, strings.Join(SupportedRuleIDs(), ", "))
	}
	if unsupportedRuleIDs[id] {
		return fmt.Errorf("rule %q requires resource kinds this integration does not load (NetworkPolicy, SecurityContextConstraints, Service, Ingress, Route or Secret objects); supported rules: %s", id, strings.Join(SupportedRuleIDs(), ", "))
	}
	return nil
}

// ValidateCustom performs kube-chainsaw-specific validation on top of the generated validation.
// This is called by the template's Instantiate function.
func (p *Params) ValidateCustom() error {
	known := make(map[string]bool)
	for _, id := range analyzer.KnownRuleIDs() {
		known[id] = true
	}

	for _, id := range p.Rules {
		if err := validateRuleID(id, known, "rule"); err != nil {
			return err
		}
	}
	for _, id := range p.ExcludeRules {
		if err := validateRuleID(id, known, "excludeRules"); err != nil {
			return err
		}
	}

	// Check for overlap
	rulesSet := make(map[string]bool)
	for _, id := range p.Rules {
		rulesSet[id] = true
	}
	for _, id := range p.ExcludeRules {
		if rulesSet[id] {
			return fmt.Errorf("rule %q appears in both rules and excludeRules", id)
		}
	}

	if p.MinSeverity != "" {
		if _, err := ParseSeverityLevel(p.MinSeverity); err != nil {
			return err
		}
	}

	return nil
}

// ParseSeverityLevel converts a severity string to a numeric level (0-3).
// Exported for use by the template's filterBySeverity function.
func ParseSeverityLevel(s string) (int, error) {
	switch strings.ToLower(s) {
	case "info", "note":
		return 0, nil
	case "warning":
		return 1, nil
	case "high", "error":
		return 2, nil
	case "critical":
		return 3, nil
	default:
		return 0, fmt.Errorf("invalid minSeverity %q; valid values: info, warning, high, critical (or SARIF: note, warning, error)", s)
	}
}
