package kubelinter.template.envvar

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	some envVar in container.env
	name := data.envvar.name
	value := data.envvar.value
	regex.match(name, envVar.name)
	regex.match(value, envVar.value)
	msg := sprintf("environment variable %s in container %q found", [envVar.name, container.name])
}