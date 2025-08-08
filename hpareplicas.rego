package kubelinter.template.hpareplicas

import data.kubelinter.objectkinds.is_horizontalpodautoscaler
import future.keywords.in

deny contains msg if {
	is_horizontalpodautoscaler
	minReplicas := data.hpareplicas.minReplicas
	replicas := get_hpa_min_replicas()
	replicas < minReplicas
	replicaText := get_replica_text(replicas)
	msg := sprintf("object has %d %s but minimum required replicas is %d", [replicas, replicaText, minReplicas])
}

get_hpa_min_replicas() := replicas if {
	replicas := input.spec.minReplicas
}

get_hpa_min_replicas() := 1 if {
	not input.spec.minReplicas
}

get_replica_text(replicas) := "replicas" if {
	replicas > 1
}

get_replica_text(replicas) := "replica" if {
	replicas == 1
}