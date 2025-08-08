package kubelinter.template.danglingingress

import kubelinter.objectkinds.is_ingress

deny contains msg if {
	is_ingress
	some serviceRef in get_ingress_service_references()
	not service_exists(serviceRef)
	msg := sprintf("no service found matching ingress label (%v), port %s", [serviceRef.name, serviceRef.port])
}

get_ingress_service_references() := refs {
	# Get default backend
	input.spec.defaultBackend
	service := input.spec.defaultBackend.service
	refs := [{"name": service.name, "port": service.port.name}]
}

get_ingress_service_references() := refs {
	# Get service references from rules
	some rule in input.spec.rules
	rule.http
	some path in rule.http.paths
	service := path.backend.service
	refs := [{"name": service.name, "port": service.port.name}]
}

service_exists(serviceRef) {
	some service in data.objects
	service.kind == "Service"
	service.metadata.namespace == input.metadata.namespace
	service.metadata.name == serviceRef.name
	some port in service.spec.ports
	port.name == serviceRef.port
}