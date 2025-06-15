package extract

import (
	"golang.stackrox.io/kube-linter/pkg/k8sutil"
	batchV1 "k8s.io/api/batch/v1"
)

// JobSpec extracts a job template spec from Job or CronJob objects
func JobSpec(obj k8sutil.Object) (batchV1.JobSpec, string, bool) {
	switch obj := obj.(type) {
	case *batchV1.Job:
		return obj.Spec, "Job", true
	case *batchV1.CronJob:
		return obj.Spec.JobTemplate.Spec, "CronJob", true
	default:
		return batchV1.JobSpec{}, "", false
	}
}
