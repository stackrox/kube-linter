package kubelinter.template.readinessport

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	container.readinessProbe
	not has_readiness_port(container)
	msg := sprintf("container %q has readiness probe but no port specified", [container.name])
}

has_readiness_port(container) if {
	some port in container.ports
	port.name == container.readinessProbe.httpGet.port
}

has_readiness_port(container) if {
	container.readinessProbe.httpGet.port
	not container.readinessProbe.httpGet.port.name
}
