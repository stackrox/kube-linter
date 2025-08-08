package kubelinter.template.hostpid

import kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	input.spec.template.spec.hostPID == true
	msg := "object shares the host's process namespace (via hostPID=true)."
}