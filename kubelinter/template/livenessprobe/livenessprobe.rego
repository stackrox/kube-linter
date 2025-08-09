package kubelinter.template.livenessprobe

import data.kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	not container.livenessProbe
	msg := sprintf("container %q does not specify a liveness probe", [container.name])
}
