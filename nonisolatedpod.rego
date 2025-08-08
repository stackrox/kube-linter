package kubelinter.template.nonisolatedpod

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in
import future.keywords.every

deny contains msg if {
	is_deployment_like
	not has_network_policy()
	msg := "object has no network policy specified"
}

has_network_policy() if {
	some np in data.objects
	np.kind == "NetworkPolicy"
	np.metadata.namespace == input.metadata.namespace
	selector_matches_pod(np.spec.podSelector, input.spec.template.metadata.labels)
}

selector_matches_pod(selector, labels) if {
	every key, value in selector.matchLabels {
		labels[key] == value
	}
}