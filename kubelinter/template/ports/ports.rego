package kubelinter.template.ports

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	some port in container.ports
	port_number := data.ports.port
	port.containerPort == port_number
	protocol := data.ports.protocol
	port_protocol := port_protocol_value(port)
	regex.match(protocol, port_protocol)
	msg := sprintf("port %d and protocol %s in container %q found", [port.containerPort, port_protocol, container.name])
}

default port_protocol_value(_) := "TCP"

port_protocol_value(port) := port.protocol
