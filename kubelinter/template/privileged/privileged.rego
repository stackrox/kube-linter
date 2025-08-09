package kubelinter.template.privileged

import data.kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	container.securityContext.privileged == true
	msg := sprintf("container %q is privileged", [container.name])
}
