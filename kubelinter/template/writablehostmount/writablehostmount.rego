package kubelinter.template.writablehostmount

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny contains msg if {
	is_deployment_like
	some volume in input.spec.template.spec.volumes
	volume.hostPath
	some container in input.spec.template.spec.containers
	some mount in container.volumeMounts
	mount.name == volume.name
	not mount.readOnly
	msg := sprintf("container %s mounts path %s on the host as writable", [container.name, volume.hostPath.path])
}
