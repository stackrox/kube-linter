package kubelinter.template.runasnonroot

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	not is_run_as_non_root(container)
	msg := sprintf("container %q is not set to runAsNonRoot", [container.name])
}

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	is_run_as_non_root_conflict(container)
	run_as_user_val := run_as_user(container)
	msg := sprintf("container %q is set to runAsNonRoot, but runAsUser set to %d", [container.name, run_as_user_val])
}

is_run_as_non_root(container) if {
	container.securityContext.runAsNonRoot == true
}

is_run_as_non_root_conflict(container) if {
	container.securityContext.runAsNonRoot == true
	container.securityContext.runAsUser == 0
}

run_as_user(container) := container.securityContext.runAsUser
