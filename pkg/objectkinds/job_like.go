package objectkinds

import (
	batchV1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	JobLike = "JobLike"
)

var (
	jobGVK = batchV1.SchemeGroupVersion.WithKind("Job")
)
var (
	cronJobGVK = batchV1.SchemeGroupVersion.WithKind("CronJob")
)

func init() {
	RegisterObjectKind(JobLike, MatcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == jobGVK || gvk == cronJobGVK
	}))
}
