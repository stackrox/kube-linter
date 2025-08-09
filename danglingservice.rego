package kubelinter.template.danglingservice

import data.kubelinter.objectkinds.is_service
import future.keywords.every
import future.keywords.in

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
	not has_matching_pods
	msg := sprintf("no pods found matching service labels (%v)", [input.spec.selector])
}

has_matching_pods if {
	some deployment in data.objects
	deployment.kind == "Deployment"
	deployment.metadata.namespace == input.metadata.namespace
	labels_match_selector(input.spec.selector[_], deployment.spec.template.metadata.labels)
}

labels_match_selector(s, labels) if labels[s.key] == s.value
