package kubelinter.template.readonlyrootfs

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	not container.securityContext.readOnlyRootFilesystem
	msg := sprintf("container %q is not set to readOnlyRootFilesystem=true", [container.name])
}
