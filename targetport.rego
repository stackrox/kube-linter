package kubelinter.template.targetport

import kubelinter.objectkinds.is_deployment_like
import kubelinter.objectkinds.is_service

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	some port in container.ports
	port.name != ""
	not is_valid_port_name(port.name)
	msg := sprintf("port name %q in container %q is invalid", [port.name, container.name])
}

deny contains msg if {
	is_service
	some port in input.spec.ports
	port.targetPort
	port.targetPort.type == "String"
	not is_valid_port_name(port.targetPort.strVal)
	msg := sprintf("port targetPort %q in service %q is invalid", [port.targetPort.strVal, input.metadata.name])
}

deny contains msg if {
	is_service
	some port in input.spec.ports
	port.targetPort
	port.targetPort.type == "Int"
	port.targetPort.intVal != 0
	not is_valid_port_number(port.targetPort.intVal)
	msg := sprintf("port targetPort %q in service %q is invalid", [port.targetPort.intVal, input.metadata.name])
}

# Simplified validation functions - in practice these would need more complex logic
is_valid_port_name(name) {
	# Basic validation: alphanumeric and hyphens only, max 15 chars
	regex.match("^[a-z0-9-]+$", name)
	count(name) <= 15
}

is_valid_port_number(port) {
	port > 0
	port <= 65535
}