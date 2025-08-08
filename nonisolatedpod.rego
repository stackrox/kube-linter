package kubelinter.template.nonisolatedpod

import kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	not is_pod_isolated_by_networkpolicy()
	msg := "pods created by this object are non-isolated"
}

is_pod_isolated_by_networkpolicy() {
	some networkpolicy in data.objects
	networkpolicy.kind == "NetworkPolicy"
	networkpolicy.metadata.namespace == input.metadata.namespace
	networkpolicy.spec.podSelector
	labels_match_selector(networkpolicy.spec.podSelector, input.spec.template.metadata.labels)
}

labels_match_selector(selector, labels) {
	# Simplified label matching - in practice this would need more complex logic
	every key, value in selector.matchLabels {
		labels[key] == value
	}
}