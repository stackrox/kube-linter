package kubelinter.template.targetport

import data.kubelinter.objectkinds.is_service
import future.keywords.in

deny contains msg if {
	is_service
	some port in input.spec.ports
	port.targetPort
	not is_valid_target_port(port.targetPort)
	msg := sprintf("service has invalid target port %v", [port.targetPort])
}

deny contains msg if {
	is_service
	some port in input.spec.ports
	port.targetPort
	not port_exists_in_pods(port.targetPort)
	msg := sprintf("target port %v not found in any pod", [port.targetPort])
}

deny contains msg if {
	is_service
	some port in input.spec.ports
	not port.targetPort
	msg := sprintf("service port %v has no target port specified", [port.port])
}

is_valid_target_port(targetPort) if {
	regex.match("^[0-9]+$", targetPort)
}

is_valid_target_port(targetPort) if {
	regex.match("^[a-zA-Z0-9-]+$", targetPort)
}

port_exists_in_pods(targetPort) if {
	some pod in data.objects
	pod.kind == "Pod"
	pod.metadata.namespace == input.metadata.namespace
	some containerPort in pod.spec.containers[0].ports
	containerPort.containerPort == targetPort
}