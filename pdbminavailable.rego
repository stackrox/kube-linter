package kubelinter.template.pdbminavailable

import kubelinter.objectkinds.is_poddisruptionbudget

deny contains msg if {
	is_poddisruptionbudget
	input.spec.minAvailable
	value := input.spec.minAvailable
	is_percentage_100(value)
	msg := "PDB has minimum available replicas set to 100 percent of replicas"
}

deny contains msg if {
	is_poddisruptionbudget
	input.spec.minAvailable
	value := input.spec.minAvailable
	replicas := get_matching_deployment_replicas()
	replicas <= value
	msg := sprintf("The current number of replicas for deployment is equal to or lower than the minimum number of replicas specified by its PDB.", [])
}

is_percentage_100(value) {
	value == "100%"
}

get_matching_deployment_replicas() := replicas {
	some deployment in data.objects
	deployment.kind == "Deployment"
	deployment.metadata.namespace == input.metadata.namespace
	labels_match_selector(input.spec.selector, deployment.spec.selector)
	replicas := deployment.spec.replicas
}

get_matching_deployment_replicas() := 1 {
	some deployment in data.objects
	deployment.kind == "Deployment"
	deployment.metadata.namespace == input.metadata.namespace
	labels_match_selector(input.spec.selector, deployment.spec.selector)
	not deployment.spec.replicas
}

labels_match_selector(selector, deploymentSelector) {
	# Simplified label matching - in practice this would need more complex logic
	every key, value in selector.matchLabels {
		deploymentSelector.matchLabels[key] == value
	}
}