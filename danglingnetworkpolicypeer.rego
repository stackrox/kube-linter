package kubelinter.template.danglingnetworkpolicypeer

import kubelinter.objectkinds.is_networkpolicy

deny contains msg if {
	is_networkpolicy
	some rule in input.spec.ingress
	some peer in rule.from
	peer.podSelector
	not peer.namespaceSelector
	not has_matching_pods_for_peer(peer.podSelector)
	msg := sprintf("no pods found matching networkpolicy rule's podSelector labels (%v)", [peer.podSelector])
}

deny contains msg if {
	is_networkpolicy
	some rule in input.spec.egress
	some peer in rule.to
	peer.podSelector
	not peer.namespaceSelector
	not has_matching_pods_for_peer(peer.podSelector)
	msg := sprintf("no pods found matching networkpolicy rule's podSelector labels (%v)", [peer.podSelector])
}

has_matching_pods_for_peer(podSelector) {
	some deployment in data.objects
	deployment.kind == "Deployment"
	deployment.metadata.namespace == input.metadata.namespace
	labels_match_selector(podSelector, deployment.spec.template.metadata.labels)
}

labels_match_selector(selector, labels) {
	# Simplified label matching - in practice this would need more complex logic
	every key, value in selector.matchLabels {
		labels[key] == value
	}
}