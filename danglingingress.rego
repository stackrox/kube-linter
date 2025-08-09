package kubelinter.template.danglingingress

import data.kubelinter.objectkinds.is_ingress
import future.keywords.in

deny contains msg if {
	is_ingress
	some serviceRef in get_ingress_service_references()
	not service_exists(serviceRef)
	msg := sprintf("no services found matching the ingress's backend service name (%s)", [serviceRef])
}

get_ingress_service_references := refs if {
	# Get default backend
	input.spec.defaultBackend
	service := input.spec.defaultBackend.service
	refs := [{"name": service.name, "port": service.port.name}]
}

get_ingress_service_references := refs if {
	# Get service references from rules
	some rule in input.spec.rules
	rule.http
	some path in rule.http.paths
	service := path.backend.service
	refs := [{"name": service.name, "port": service.port.name}]
}

service_exists(serviceRef) if {
	some service in data.objects
	service.kind == "Service"
	service.metadata.namespace == input.metadata.namespace
	service.metadata.name == serviceRef.name
	some port in service.spec.ports
	port.name == serviceRef.port
}
