package kubelinter.template.replicas

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny contains msg if {
	is_deployment_like
	replicas := input.spec.replicas
	replicas < data.replicas.minReplicas
	msg := sprintf("object has %d replicas but minimum required is %d", [replicas, data.replicas.minReplicas])
}