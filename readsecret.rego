package kubelinter.template.readsecret

import kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	some envVar in container.env
	envVar.valueFrom.secretKeyRef
	msg := sprintf("environment variable %q in container %q uses SecretKeyRef", [envVar.name, container.name])
}