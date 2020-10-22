package util

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/diagnostic"
	"golang.stackrox.io/kube-linter/internal/extract"
	"golang.stackrox.io/kube-linter/internal/k8sutil"
	"golang.stackrox.io/kube-linter/internal/lintcontext"
	"golang.stackrox.io/kube-linter/internal/matcher"
	"golang.stackrox.io/kube-linter/internal/stringutils"
)

// ConstructRequiredMapMatcher constructs a check function that requires that a k-v pair is present in the map.
func ConstructRequiredMapMatcher(key, value, fieldType string) (check.Func, error) {
	keyMatcher, err := matcher.ForString(key)
	if err != nil {
		return nil, errors.Wrap(err, "invalid key")
	}
	valueMatcher, err := matcher.ForString(value)
	if err != nil {
		return nil, errors.Wrap(err, "invalid value")
	}

	var extractFunc func(object k8sutil.Object) map[string]string
	switch fieldType {
	case "label":
		extractFunc = extract.Labels
	case "annotation":
		extractFunc = extract.Annotations
	default:
		return nil, errors.Errorf("unknown fieldType %q", fieldType)
	}

	return func(_ *lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
		fields := extractFunc(object.K8sObject)
		for k, v := range fields {
			if keyMatcher(k) && valueMatcher(v) {
				return nil
			}
		}
		return []diagnostic.Diagnostic{{
			Message: fmt.Sprintf("no %s matching \"%s=%s\" found", fieldType, key, stringutils.OrDefault(value, "<any>")),
		}}
	}, nil
}
