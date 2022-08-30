package objectkinds

import (
	"fmt"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	allObjectKinds = make(map[string]Matcher)
)

// RegisterObjectKind allows a matcher function to be registered for a given object kind
func RegisterObjectKind(name string, objectKind Matcher) {
	if _, ok := allObjectKinds[name]; ok {
		panic(fmt.Sprintf("duplicate object kind: %v", name))
	}
	allObjectKinds[name] = objectKind
}

// AllObjectKinds will return all the object kind names that are registered
func AllObjectKinds() []string {
	kinds := make([]string, 0, len(allObjectKinds))

	for k := range allObjectKinds {
		kinds = append(kinds, k)
	}
	return kinds
}

type orMatcher []Matcher

func (o orMatcher) Matches(gvk schema.GroupVersionKind) bool {
	for _, m := range o {
		if m.Matches(gvk) {
			return true
		}
	}
	return false
}

// ConstructMatcher constructs a matcher that matches objects that fall
// into one of the given object kinds.
func ConstructMatcher(objectKinds ...string) (Matcher, error) {
	var matchers []Matcher
	for _, obj := range objectKinds {
		matcher := allObjectKinds[obj]
		if matcher == nil {
			return nil, errors.Errorf("unknown object kind: %v", obj)
		}
		matchers = append(matchers, matcher)
	}
	return orMatcher(matchers), nil
}
