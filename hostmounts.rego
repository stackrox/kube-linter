package kubelinter.template.hostmounts

import kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	some volume in input.spec.template.spec.volumes
	volume.hostPath
	some dirPattern in data.hostmounts.dirs
	regex.match(dirPattern, volume.hostPath.path)
	some container in input.spec.template.spec.containers
	some mount in container.volumeMounts
	mount.name == volume.name
	msg := sprintf("host system directory %q is mounted on container %q", [volume.hostPath.path, container.name])
}