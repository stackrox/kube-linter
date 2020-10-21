package ignore

import (
	"golang.stackrox.io/kube-linter/internal/stringutils"
)

const (
	// AnnotationKeyPrefix is the prefix for annotations for kube-linter check ignores.
	AnnotationKeyPrefix = "kube-linter.io/ignore-"

	// All is a special term used to indicate that _all_ checks are to be ignored for the given object.
	All = "all"
)

// ObjectForCheck returns whether to ignore the given object for the passed check name.
func ObjectForCheck(annotations map[string]string, checkName string) bool {
	for k := range annotations {
		key := k
		if stringutils.ConsumePrefix(&key, AnnotationKeyPrefix) {
			if key == All || key == checkName {
				return true
			}
		}
	}
	return false
}
