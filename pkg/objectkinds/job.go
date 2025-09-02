package objectkinds

import (
	batchV1 "k8s.io/api/batch/v1"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// Job represents Kubernetes Job objects.
	Job = "Job"
)

var (
	jobGVK = batchV1.SchemeGroupVersion.WithKind("Job")
)

func init() {
	RegisterObjectKind(Job, MatcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == jobGVK
	}))
}
