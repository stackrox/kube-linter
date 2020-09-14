package objectkinds

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// A Matcher selects a certain subset of GVKs.
type Matcher interface {
	Matches(gvk schema.GroupVersionKind) bool
}

type matcherFunc func(gvk schema.GroupVersionKind) bool

func (f matcherFunc) Matches(gvk schema.GroupVersionKind) bool {
	return f(gvk)
}
