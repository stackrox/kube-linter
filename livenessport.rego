package kubelinter.template.livenessport

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	container.livenessProbe
	not has_liveness_port(container)
	msg := sprintf("container %q has liveness probe but no port specified", [container.name])
}

has_liveness_port(container) if {
	some port in container.ports
	port.name == container.livenessProbe.httpGet.port
}

has_liveness_port(container) if {
	container.livenessProbe.httpGet.port
	not container.livenessProbe.httpGet.port.name
}
