package kubelinter.template.restartpolicy

import data.kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	restart_policy := input.spec.template.spec.restartPolicy
	not is_accepted_restart_policy(restart_policy)
	msg := sprintf(
		"object has a restart policy defined with '%s' but the only accepted restart policies are 'Always' and 'OnFailure'",
		[restart_policy],
	)
}

is_accepted_restart_policy("Always")

is_accepted_restart_policy("OnFailure")
