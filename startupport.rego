package kubelinter.template.startupport

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	container.startupProbe
	not has_startup_port(container)
	msg := sprintf("container %q has startup probe but no port specified", [container.name])
}

has_startup_port(container) if {
	some port in container.ports
	port.name == container.startupProbe.httpGet.port
}

has_startup_port(container) if {
	container.startupProbe.httpGet.port
	not container.startupProbe.httpGet.port.name
}
