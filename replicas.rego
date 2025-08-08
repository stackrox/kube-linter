package kubelinter.template.replicas

import kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	minReplicas := data.replicas.minReplicas
	replicas := get_replicas()
	replicas < minReplicas
	replicaText := get_replica_text(replicas)
	msg := sprintf("object has %d %s but minimum required replicas is %d", [replicas, replicaText, minReplicas])
}

get_replicas() := replicas {
	# Check for explicit replicas field
	replicas := input.spec.replicas
}

get_replicas() := 1 {
	# Default to 1 if not specified
	not input.spec.replicas
}

get_replica_text(replicas) := "replicas" {
	replicas > 1
}

get_replica_text(replicas) := "replica" {
	replicas == 1
}