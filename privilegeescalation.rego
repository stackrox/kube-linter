package kubelinter.template.privilegeescalation

import kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	container.securityContext.allowPrivilegeEscalation == true
	msg := sprintf("container %q has AllowPrivilegeEscalation set to true.", [container.name])
}

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	container.securityContext.privileged == true
	msg := sprintf("container %q is Privileged hence allows privilege escalation.", [container.name])
}

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	some capability in container.securityContext.capabilities.add
	capability == "SYS_ADMIN"
	msg := sprintf("container %q has SYS_ADMIN capability hence allows privilege escalation.", [container.name])
}