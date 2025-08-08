package kubelinter.template.duplicatenvvar

import data.kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	some envVar in container.env
	count(container.env, envVar.name) > 1
	msg := sprintf("Duplicate environment variable %s in container %q found", [envVar.name, container.name])
}