package objectkinds

import (
	batchV1 "k8s.io/api/batch/v1"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// CronJob represents Kubernetes CronJob objects.
	CronJob = "CronJob"
)

var (
	cronJobGVK = batchV1.SchemeGroupVersion.WithKind("CronJob")
)

func init() {
	RegisterObjectKind(CronJob, MatcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == cronJobGVK
	}))
}
