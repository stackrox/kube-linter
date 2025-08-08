package kubelinter.template.containercapabilities

import kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	container.securityContext.capabilities
	forbiddenCaps := data.containercapabilities.forbiddenCapabilities
	exceptions := data.containercapabilities.exceptions
	check_forbidden_capabilities(container, forbiddenCaps, exceptions)
	msg := sprintf("container %q has ADD capability: %q, which matched with the forbidden capability for containers", [container.name, get_added_capability(container)])
}

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	container.securityContext.capabilities
	forbiddenCaps := data.containercapabilities.forbiddenCapabilities
	check_missing_drop_capabilities(container, forbiddenCaps)
	msg := sprintf("container %q has DROP capabilities: %q, but does not drop capability %q which is required", [container.name, container.securityContext.capabilities.drop, get_missing_drop_capability(container, forbiddenCaps)])
}

check_forbidden_capabilities(container, forbiddenCaps, exceptions) {
	some forbiddenCap in forbiddenCaps
	forbiddenCap == "all"
	some addedCap in container.securityContext.capabilities.add
	not is_exception(addedCap, exceptions)
}

check_forbidden_capabilities(container, forbiddenCaps, exceptions) {
	some forbiddenCap in forbiddenCaps
	forbiddenCap != "all"
	some addedCap in container.securityContext.capabilities.add
	regex.match(forbiddenCap, addedCap)
}

check_missing_drop_capabilities(container, forbiddenCaps) {
	some forbiddenCap in forbiddenCaps
	forbiddenCap == "all"
	not has_all_drop(container)
}

check_missing_drop_capabilities(container, forbiddenCaps) {
	some forbiddenCap in forbiddenCaps
	forbiddenCap != "all"
	not has_drop_capability(container, forbiddenCap)
}

is_exception(cap, exceptions) {
	some exception in exceptions
	regex.match(exception, cap)
}

has_all_drop(container) {
	some dropCap in container.securityContext.capabilities.drop
	dropCap == "all"
}

has_drop_capability(container, forbiddenCap) {
	some dropCap in container.securityContext.capabilities.drop
	regex.match(forbiddenCap, dropCap)
}

has_drop_capability(container, forbiddenCap) {
	some dropCap in container.securityContext.capabilities.drop
	dropCap == "all"
}

get_added_capability(container) := cap {
	some cap in container.securityContext.capabilities.add
}

get_missing_drop_capability(container, forbiddenCaps) := cap {
	some cap in forbiddenCaps
	not has_drop_capability(container, cap)
}