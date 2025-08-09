package kubelinter.template.pdbminavailable

import data.kubelinter.objectkinds.is_poddisruptionbudget
import future.keywords.every
import future.keywords.in

deny contains msg if {
	is_poddisruptionbudget
	minAvailable := input.spec.minAvailable
	not is_valid_min_available(minAvailable)
	msg := sprintf("minAvailable %v is not valid", [minAvailable])
}

deny contains msg if {
	is_poddisruptionbudget
	minAvailable := input.spec.minAvailable
	not has_matching_deployments(minAvailable)
	msg := sprintf("no deployments found matching PDB minAvailable %v", [minAvailable])
}

is_valid_min_available(minAvailable) if {
	minAvailable > 0
}

is_valid_min_available(minAvailable) if {
	minAvailable == "100%"
}

has_matching_deployments(minAvailable) if {
	some deployment in data.objects
	deployment.kind == "Deployment"
	deployment.metadata.namespace == input.metadata.namespace
	labels_match_selector(input.spec.selector, deployment.spec.template.metadata.labels)
}

labels_match_selector(selector, labels) if {
	every key, value in selector.matchLabels {
		labels[key] == value
	}
}
