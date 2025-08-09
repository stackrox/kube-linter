package kubelinter.template.pdbminavailable

import data.kubelinter.objectkinds.is_poddisruptionbudget
import future.keywords.every
import future.keywords.in

deny contains msg if {
	is_poddisruptionbudget
	min_available := input.spec.minAvailable
	not is_valid_min_available(min_available)
	msg := sprintf("minAvailable %v is not valid", [min_available])
}

deny contains msg if {
	is_poddisruptionbudget
	min_available := input.spec.minAvailable
	not has_matching_deployments(min_available)
	msg := sprintf("no deployments found matching PDB minAvailable %v", [min_available])
}

is_valid_min_available(min_available) if {
	min_available > 0
}

is_valid_min_available("100%")

has_matching_deployments(min_available) if {
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
