package kubelinter.template.readsecret

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in
deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	some envVar in container.env
	envVar.valueFrom.secretKeyRef
	msg := sprintf("environment variable %q in container %q uses SecretKeyRef", [envVar.name, container.name])
}