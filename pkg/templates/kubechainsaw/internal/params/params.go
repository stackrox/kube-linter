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

// ValidateCustom performs kube-chainsaw-specific validation on top of the generated validation.
// This is called by the template's Instantiate function.
func (p *Params) ValidateCustom() error {
	known := make(map[string]bool)
	for _, id := range analyzer.KnownRuleIDs() {
		known[id] = true
	}

	for _, id := range p.Rules {
		if !known[id] {
			return fmt.Errorf("unknown rule ID %q; valid IDs: %s", id, strings.Join(analyzer.KnownRuleIDs(), ", "))
		}
	}
	for _, id := range p.ExcludeRules {
		if !known[id] {
			return fmt.Errorf("unknown excludeRules ID %q; valid IDs: %s", id, strings.Join(analyzer.KnownRuleIDs(), ", "))
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
		if _, err := parseMinSeverity(p.MinSeverity); err != nil {
			return err
		}
	}

	return nil
}

func parseMinSeverity(s string) (int, error) {
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
