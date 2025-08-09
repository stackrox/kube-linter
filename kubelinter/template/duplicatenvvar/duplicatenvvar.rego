package kubelinter.template.duplicatenvvar

import data.kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	some env_var in container.env
	env_var_count := count([v | some v in container.env; v.name == env_var.name])
	env_var_count > 1
	msg := sprintf("Duplicate environment variable %s in container %q found", [env_var.name, container.name])
}
