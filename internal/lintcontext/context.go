package lintcontext

import (
	"golang.stackrox.io/kube-linter/internal/k8sutil"
)

// ObjectMetadata is metadata about an object.
type ObjectMetadata struct {
	FilePath string
	Raw      []byte
}

// An Object references an object that is loaded from a YAML file.
type Object struct {
	Metadata  ObjectMetadata
	K8sObject k8sutil.Object
}

// An InvalidObject represents something that couldn't be loaded from a YAML file.
type InvalidObject struct {
	Metadata ObjectMetadata
	LoadErr  error
}

// A LintContext represents the context for a lint run.
type LintContext struct {
	objects        []Object
	invalidObjects []InvalidObject
}

// New returns a ready-to-use, empty, lint context.
func New() *LintContext {
	return &LintContext{}
}
