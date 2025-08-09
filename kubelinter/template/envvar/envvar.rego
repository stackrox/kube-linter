package kubelinter.template.envvar

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	some env_var in container.env
	name := data.envvar.name
	regex.match(name, env_var.name)
	value := data.envvar.value
	regex.match(value, env_var.value)
	msg := sprintf("environment variable %s in container %q found", [env_var.name, container.name])
}
