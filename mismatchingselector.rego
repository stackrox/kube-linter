package kubelinter.template.mismatchingselector

import kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	not is_job_or_cronjob()
	not has_selector()
	msg := "object has no selector specified"
}

deny contains msg if {
	is_deployment_like
	has_selector()
	not selector_matches_pod_labels()
	msg := sprintf("labels in pod spec (%v) do not match labels in selector (%v)", [input.spec.template.metadata.labels, input.spec.selector])
}

is_job_or_cronjob() {
	input.kind == "Job"
}

is_job_or_cronjob() {
	input.kind == "CronJob"
}

has_selector() {
	input.spec.selector
	count(input.spec.selector.matchLabels) > 0
}

has_selector() {
	input.spec.selector
	count(input.spec.selector.matchExpressions) > 0
}

selector_matches_pod_labels() {
	# Simplified check - in practice this would need more complex label matching logic
	every key, value in input.spec.selector.matchLabels {
		input.spec.template.metadata.labels[key] == value
	}
}