package objectkinds

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// CronJob represents Kubernetes CronJob objects.
	CronJob = "CronJob"
)

func init() {
	RegisterObjectKind(CronJob, MatcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == cronJobGVK
	}))
}
