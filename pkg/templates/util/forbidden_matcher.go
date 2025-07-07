package util

import (
	"fmt"

	"golang.stackrox.io/kube-linter/internal/stringutils"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/k8sutil"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/matcher"
)

// ConstructForbiddenMapMatcher constructs a check function that requires that a k-v pair is NOT present in the map.
func ConstructForbiddenMapMatcher(key, value, fieldType string) (check.Func, error) {
	keyMatcher, err := matcher.ForString(key)
	if err != nil {
		return nil, fmt.Errorf("invalid key: %w", err)
	}
	valueMatcher, err := matcher.ForString(value)
	if err != nil {
		return nil, fmt.Errorf("invalid value: %w", err)
	}

	var extractFunc func(object k8sutil.Object) map[string]string
	switch fieldType {
	case "label":
		extractFunc = extract.Labels
	case "annotation":
		extractFunc = extract.Annotations
	default:
		return nil, fmt.Errorf("unknown fieldType %q", fieldType)
	}

	return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
		fields := extractFunc(object.K8sObject)
		for k, v := range fields {
			if keyMatcher(k) && valueMatcher(v) {
				return []diagnostic.Diagnostic{{
					Message: fmt.Sprintf("%s matching \"%s=%s\" found", fieldType, key, stringutils.OrDefault(value, "<any>")),
				}}
			}
		}
		return nil
	}, nil
}
