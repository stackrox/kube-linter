package objectkinds

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// A Matcher selects a certain subset of GVKs.
type Matcher interface {
	Matches(gvk schema.GroupVersionKind) bool
}

// MatcherFunc takes in a GVK and decides if it matches an object kind
type MatcherFunc func(gvk schema.GroupVersionKind) bool

func (f MatcherFunc) Matches(gvk schema.GroupVersionKind) bool {
	return f(gvk)
}
