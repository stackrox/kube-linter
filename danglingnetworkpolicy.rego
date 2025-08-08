package kubelinter.template.danglingnetworkpolicy

import kubelinter.objectkinds.is_networkpolicy

deny contains msg if {
	is_networkpolicy
	selector := input.spec.podSelector
	not has_matching_pods(selector)
	msg := sprintf("no pods found matching networkpolicy's podSelector labels (%v)", [selector])
}

has_matching_pods(selector) {
	some deployment in data.objects
	deployment.kind == "Deployment"
	deployment.metadata.namespace == input.metadata.namespace
	labels_match_selector(selector, deployment.spec.template.metadata.labels)
}

labels_match_selector(selector, labels) {
	# Simplified label matching - in practice this would need more complex logic
	every key, value in selector.matchLabels {
		labels[key] == value
	}
}