package kubelinter.template.runasnonroot

import kubelinter.objectkinds.is_deployment_like

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
	runAsUser := get_run_as_user(container)
	msg := sprintf("container %q is set to runAsNonRoot, but runAsUser set to %d", [container.name, runAsUser])
}

is_run_as_non_root(container) {
	# Check if runAsNonRoot is explicitly set to true
	container.securityContext.runAsNonRoot == true
}

is_run_as_non_root(container) {
	# Check if runAsNonRoot is set at pod level
	input.spec.template.spec.securityContext.runAsNonRoot == true
}

is_run_as_non_root(container) {
	# Check if runAsUser is explicitly set to non-root (> 0)
	runAsUser := get_run_as_user(container)
	runAsUser > 0
}

is_run_as_non_root_conflict(container) {
	# Check if runAsNonRoot is set but runAsUser is set to 0
	(is_run_as_non_root(container) | input.spec.template.spec.securityContext.runAsNonRoot == true)
	runAsUser := get_run_as_user(container)
	runAsUser == 0
}

get_run_as_user(container) := user {
	# Check container level first
	user := container.securityContext.runAsUser
}

get_run_as_user(container) := user {
	# Fall back to pod level
	user := input.spec.template.spec.securityContext.runAsUser
}

get_run_as_user(container) := 0 {
	# Default to 0 if not set
	not container.securityContext.runAsUser
	not input.spec.template.spec.securityContext.runAsUser
}