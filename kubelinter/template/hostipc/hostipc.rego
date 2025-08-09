package kubelinter.template.hostipc

import data.kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	input.spec.template.spec.hostIPC == true
	msg := "resource shares host's IPC namespace (via hostIPC=true)."
}
