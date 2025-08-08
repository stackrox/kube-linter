package kubelinter.template.imagepullpolicy

import kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	some forbiddenPolicy in data.imagepullpolicy.forbiddenPolicies
	container.imagePullPolicy == forbiddenPolicy
	msg := sprintf("container %q has imagePullPolicy set to %s", [container.name, container.imagePullPolicy])
}