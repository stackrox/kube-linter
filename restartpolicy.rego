package kubelinter.template.restartpolicy

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny contains msg if {
	is_deployment_like
	restartPolicy := input.spec.template.spec.restartPolicy
	restartPolicy != "Always"
	msg := sprintf("restart policy %q is not allowed", [restartPolicy])
}