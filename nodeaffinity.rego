package kubelinter.template.nodeaffinity

import kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	not input.spec.template.spec.affinity
	msg := "object does not define any node affinity rules."
}

deny contains msg if {
	is_deployment_like
	input.spec.template.spec.affinity
	not input.spec.template.spec.affinity.nodeAffinity
	msg := "object does not define any node affinity rules."
}