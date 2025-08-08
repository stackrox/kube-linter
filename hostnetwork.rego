package kubelinter.template.hostnetwork

import kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	input.spec.template.spec.hostNetwork == true
	msg := "resource shares host's network namespace (via hostNetwork=true)."
}