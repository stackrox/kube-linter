package kubelinter.template.imagepullpolicy

import data.kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	some forbidden_policy in data.imagepullpolicy.forbiddenPolicies
	container.imagePullPolicy == forbidden_policy
	msg := sprintf("container %q has imagePullPolicy set to %s", [container.name, container.imagePullPolicy])
}
