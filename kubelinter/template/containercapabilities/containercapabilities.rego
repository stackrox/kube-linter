package kubelinter.template.containercapabilities

import data.kubelinter.objectkinds
import future.keywords.in

deny contains msg if {
	objectkinds.is_deployment_like
	some container in input.spec.template.spec.containers
	container.securityContext.capabilities
	forbidden_caps := data.containercapabilities.forbidden_capabilities
	exceptions := data.containercapabilities.exceptions
	check_forbidden_capabilities(container, forbidden_caps, exceptions)
	msg := sprintf(
		"container %q has ADD capability: %q, which matched with the forbidden capability for containers",
		[container.name, added_capability(container)],
	)
}

deny contains msg if {
	objectkinds.is_deployment_like
	some container in input.spec.template.spec.containers
	container.securityContext.capabilities
	forbidden_caps := data.containercapabilities.forbidden_capabilities
	check_missing_drop_capabilities(container, forbidden_caps)
	msg := sprintf(
		"container %q has DROP capabilities: %q, but does not drop capability %q which is required",
		[container.name, container.securityContext.capabilities.drop, missing_drop_capability(container, forbidden_caps)],
	)
}

check_forbidden_capabilities(container, forbidden_caps, exceptions) if {
	some forbidden_cap in forbidden_caps
	forbidden_cap == "all"
	some added_cap in container.securityContext.capabilities.add
	not is_exception(added_cap, exceptions)
}

check_forbidden_capabilities(container, forbidden_caps, exceptions) if {
	some forbidden_cap in forbidden_caps
	forbidden_cap != "all"
	some added_cap in container.securityContext.capabilities.add
	regex.match(forbidden_cap, added_cap)
}

check_missing_drop_capabilities(container, forbidden_caps) if {
	"all" in forbidden_caps
	count([cap | some cap in container.securityContext.capabilities.drop; cap == "all"]) == 0
}

check_missing_drop_capabilities(container, forbidden_caps) if {
	some forbidden_cap in forbidden_caps
	forbidden_cap != "all"
	not has_drop_capability(container, forbidden_cap)
}

is_exception(cap, exceptions) if {
	some exception in exceptions
	regex.match(exception, cap)
}

has_all_drop(container) if {
	some drop_cap in container.securityContext.capabilities.drop
	drop_cap == "all"
}

has_drop_capability(container, forbidden_cap) if {
	some drop_cap in container.securityContext.capabilities.drop
	regex.match(forbidden_cap, drop_cap)
}

has_drop_capability(container, forbidden_cap) if {
	some drop_cap in container.securityContext.capabilities.drop
	drop_cap == "all"
}

added_capability(container) := cap if {
	some cap in container.securityContext.capabilities.add
}

missing_drop_capability(container, forbidden_caps) := cap if {
	some cap in forbidden_caps
	not has_drop_capability(container, cap)
}
