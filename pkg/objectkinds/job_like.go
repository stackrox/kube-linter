package objectkinds

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	JobLike = "JobLike"
)

func init() {
	RegisterObjectKind(JobLike, MatcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == jobGVK || gvk == cronJobGVK
	}))
}
