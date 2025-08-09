package kubelinter.template.danglingnetworkpolicy

import data.kubelinter.objectkinds.is_networkpolicy
import future.keywords.every
import future.keywords.in

deny contains msg if {
	is_networkpolicy
	selector := input.spec.podSelector
	not has_matching_pods(selector)
	msg := sprintf("no pods found matching network policy selector %v", [selector])
}

has_matching_pods(selector) if {
	some pod in data.objects
	pod.kind == "Pod"
	pod.metadata.namespace == input.metadata.namespace
	labels_match_selector(selector, pod.metadata.labels)
}

labels_match_selector(selector, labels) if {
	every key, value in selector.matchLabels {
		labels[key] == value
	}
}
