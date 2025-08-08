package kubelinter.template.ports

import kubelinter.objectkinds.is_deployment_like

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

get_port_protocol(port) := protocol {
	# Default to TCP if not specified
	protocol := port.protocol
}

get_port_protocol(port) := "TCP" {
	# Default to TCP if not specified
	not port.protocol
}