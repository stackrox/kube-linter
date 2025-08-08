package kubelinter.template.ports

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in
deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	some port in container.ports
	portNumber := data.ports.port
	protocol := data.ports.protocol
	port.containerPort == portNumber
	portProtocol := get_port_protocol(port)
	regex.match(protocol, portProtocol)
	msg := sprintf("port %d and protocol %s in container %q found", [port.containerPort, portProtocol, container.name])
}

get_port_protocol(port) :=	port.protocol
default get_port_protocol(_) = "TCP"