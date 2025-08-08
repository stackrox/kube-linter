package kubelinter.template.livenessport

import kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	container.livenessProbe
	not port_exposed_for_probe(container, container.livenessProbe)
	msg := get_probe_port_error_message(container, container.livenessProbe)
}

port_exposed_for_probe(container, probe) {
	# Get all exposed ports
	exposedPorts := get_exposed_ports(container)

	# Check HTTP probe
	probe.httpGet
	probe.httpGet.port in exposedPorts
}

port_exposed_for_probe(container, probe) {
	# Get all exposed ports
	exposedPorts := get_exposed_ports(container)

	# Check TCP probe
	probe.tcpSocket
	probe.tcpSocket.port in exposedPorts
}

port_exposed_for_probe(container, probe) {
	# Get all exposed ports
	exposedPorts := get_exposed_ports(container)

	# Check GRPC probe
	probe.grpc
	probe.grpc.port in exposedPorts
}

get_exposed_ports(container) := ports {
	ports := [port.containerPort | port := container.ports[_]]
	ports := array.concat(ports, [port.name | port := container.ports[_]])
}

get_probe_port_error_message(container, probe) := msg {
	probe.httpGet
	msg := sprintf("container %q does not expose port %s for the HTTPGet", [container.name, probe.httpGet.port])
}

get_probe_port_error_message(container, probe) := msg {
	probe.tcpSocket
	msg := sprintf("container %q does not expose port %s for the TCPSocket", [container.name, probe.tcpSocket.port])
}

get_probe_port_error_message(container, probe) := msg {
	probe.grpc
	msg := sprintf("container %q does not expose port %d for the GRPC check", [container.name, probe.grpc.port])
}