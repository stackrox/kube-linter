package kubelinter.template.jobttlsecondsafterfinished

import data.kubelinter.objectkinds.is_job_like

deny contains msg if {
	is_job_like
	input.kind == "Job"
	not input.spec.ttlSecondsAfterFinished
	msg := "Standalone Job does not specify ttlSecondsAfterFinished"
}

deny contains msg if {
	is_job_like
	input.kind == "CronJob"
	input.spec.ttlSecondsAfterFinished
	msg := "Managed Job specifies ttlSecondsAfterFinished which might conflict with successfulJobsHistoryLimit and failedJobsHistoryLimit from CronJob that have default values. Final behaviour is determined by the strictest parameter, and therefore, setting ttlSecondsAfterFinished at the job level can result with unexpected behaviour with regard to finished jobs removal"
}