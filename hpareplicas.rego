package kubelinter.template.hpareplicas

import kubelinter.objectkinds.is_horizontalpodautoscaler

deny contains msg if {
	is_horizontalpodautoscaler
	minReplicas := data.hpareplicas.minReplicas
	replicas := get_hpa_min_replicas()
	replicas < minReplicas
	replicaText := get_replica_text(replicas)
	msg := sprintf("object has %d %s but minimum required replicas is %d", [replicas, replicaText, minReplicas])
}

get_hpa_min_replicas() := replicas {
	replicas := input.spec.minReplicas
}

get_hpa_min_replicas() := 1 {
	not input.spec.minReplicas
}

get_replica_text(replicas) := "replicas" {
	replicas > 1
}

get_replica_text(replicas) := "replica" {
	replicas == 1
}