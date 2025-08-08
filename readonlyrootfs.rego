package kubelinter.template.readonlyrootfs

import kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	not has_readonly_rootfs(container)
	msg := sprintf("container %q does not have a read-only root file system", [container.name])
}

has_readonly_rootfs(container) {
	container.securityContext.readOnlyRootFilesystem == true
}