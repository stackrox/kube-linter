package kubelinter.template.danglingservice

import kubelinter.objectkinds.is_service

deny contains msg if {
	is_service
	input.spec.type != "ExternalName"
	not input.spec.selector
	msg := "service has no selector specified"
}

deny contains msg if {
	is_service
	input.spec.type != "ExternalName"
	input.spec.selector
	not has_matching_pods()
	ignoredLabels := data.danglingservice.ignoredLabels
	msg := sprintf("no pods found matching service labels (%v)", [input.spec.selector])
}

has_matching_pods() {
	some deployment in data.objects
	deployment.kind == "Deployment"
	deployment.metadata.namespace == input.metadata.namespace
	labels_match_selector(input.spec.selector, deployment.spec.template.metadata.labels)
}

labels_match_selector(selector, labels) {
	# Simplified label matching - in practice this would need more complex logic
	every key, value in selector {
		labels[key] == value
	}
}