package objectkinds

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// Any represents the ObjectKind that matches any object.
	Any = "Any"
)

func init() {
	RegisterObjectKind(Any, MatcherFunc(func(gvk schema.GroupVersionKind) bool {
		return true
	}))
}
