package objectkinds

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// Job represents Kubernetes Job objects.
	Job = "Job"
)

func init() {
	RegisterObjectKind(Job, MatcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == jobGVK
	}))
}
