package kubelinter.template.privilegedports

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	some port in container.ports
	port.containerPort > 0
	port.containerPort < 1024
	msg := sprintf("port %d is mapped in container %q.", [port.containerPort, container.name])
}
