package kubelinter.template.unsafeprocmount

import kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	container.securityContext.procMount == "Unmasked"
	msg := sprintf("container %q exposes /proc unsafely (via procMount=Unmasked).", [container.name])
}