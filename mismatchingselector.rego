package kubelinter.template.mismatchingselector

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.every
import future.keywords.in

deny contains msg if {
	is_deployment_like
	selector := input.spec.selector
	templateLabels := input.spec.template.metadata.labels
	not labels_match_selector(selector, templateLabels)
	msg := sprintf("selector %v does not match template labels %v", [selector, templateLabels])
}

deny contains msg if {
	is_deployment_like
	selector := input.spec.selector
	templateLabels := input.spec.template.metadata.labels
	not selector_matches_labels(selector, templateLabels)
	msg := sprintf("template labels %v do not match selector %v", [templateLabels, selector])
}

is_job_or_cronjob if {
	input.kind == "Job"
}

is_job_or_cronjob if {
	input.kind == "CronJob"
}

has_selector if {
	input.spec.selector
	count(input.spec.selector.matchLabels) > 0
}

has_selector if {
	input.spec.selector
	count(input.spec.selector.matchExpressions) > 0
}

selector_matches_pod_labels if {
	# Simplified check - in practice this would need more complex label matching logic
	every key, value in input.spec.selector.matchLabels {
		input.spec.template.metadata.labels[key] == value
	}
}
