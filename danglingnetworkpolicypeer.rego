package kubelinter.template.danglingnetworkpolicypeer

import data.kubelinter.objectkinds.is_networkpolicy
import future.keywords.every
import future.keywords.in

deny contains msg if {
	is_networkpolicy
	some ingress in input.spec.ingress
	some peer in ingress.from
	peer.podSelector
	not has_matching_pods_for_peer(peer.podSelector)
	msg := sprintf("no pods found matching network policy peer selector %v", [peer.podSelector])
}

deny contains msg if {
	is_networkpolicy
	some egress in input.spec.egress
	some peer in egress.to
	peer.podSelector
	not has_matching_pods_for_peer(peer.podSelector)
	msg := sprintf("no pods found matching network policy peer selector %v", [peer.podSelector])
}

has_matching_pods_for_peer(podSelector) if {
	some pod in data.objects
	pod.kind == "Pod"
	pod.metadata.namespace == input.metadata.namespace
	labels_match_selector(podSelector, pod.metadata.labels)
}

labels_match_selector(selector, labels) if {
	every key, value in selector.matchLabels {
		labels[key] == value
	}
}
