package kubelinter.template.restartpolicy

import kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	restartPolicy := input.spec.template.spec.restartPolicy
	not is_accepted_restart_policy(restartPolicy)
	msg := sprintf("object has a restart policy defined with '%s' but the only accepted restart policies are 'Always' and 'OnFailure'", [restartPolicy])
}

is_accepted_restart_policy(policy) {
	policy == "Always"
}

is_accepted_restart_policy(policy) {
	policy == "OnFailure"
}