package lintcontext

import (
	"golang.stackrox.io/kube-linter/internal/k8sutil"
)

// An ObjectWithMetadata references an object that is loaded from a YAML file.
type ObjectWithMetadata struct {
	FilePath string

	Raw []byte

	K8sObject k8sutil.Object

	LoadErr error
}

// A LintContext represents the context for a lint run.
type LintContext struct {
	Objects []ObjectWithMetadata
}

// New returns a ready-to-use, empty, lint context.
func New() *LintContext {
	return &LintContext{}
}
