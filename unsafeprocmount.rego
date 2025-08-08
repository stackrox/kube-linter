package kubelinter.template.unsafeprocmount

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in
deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	container.securityContext.procMount == "Unmasked"
	msg := sprintf("container %q exposes /proc unsafely (via procMount=Unmasked).", [container.name])
}