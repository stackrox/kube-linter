package kubelinter.template.hpareplicas

import data.kubelinter.objectkinds.is_horizontalpodautoscaler
import future.keywords.in

deny contains msg if {
	is_horizontalpodautoscaler
	min_replicas := data.hpareplicas.minReplicas
	hpa_min_replicas < min_replicas
	replica_text_val := replica_text(hpa_min_replicas)
	msg := sprintf(
		"object has %d %s but minimum required replicas is %d",
		[hpa_min_replicas, replica_text_val, min_replicas],
	)
}

hpa_min_replicas := input.spec.minReplicas

hpa_min_replicas := 1

replica_text(replicas) := "replicas" if {
	replicas > 1
}

replica_text(1) := "replica"
