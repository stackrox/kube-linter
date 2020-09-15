package matcher

import (
	"regexp"

	"golang.stackrox.io/kube-linter/internal/stringutils"
)

const (
	// NegationPrefix is the prefix used for negations.
	NegationPrefix = "!"
)

func matchAny(_ string) bool {
	return true
}

// ForString constructs a string matcher for the given value.
func ForString(value string) (func(string) bool, error) {
	if value == "" {
		return matchAny, nil
	}
	var negate bool
	if stringutils.ConsumePrefix(&value, NegationPrefix) {
		negate = true
	}
	re, err := regexp.Compile(value)
	if err != nil {
		return nil, err
	}
	return func(s string) bool {
		matched := re.MatchString(s)
		return matched != negate
	}, nil
}
