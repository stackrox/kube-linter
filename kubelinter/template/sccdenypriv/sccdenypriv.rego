package kubelinter.template.sccdenypriv

import data.kubelinter.objectkinds.is_securitycontextconstraints

deny contains msg if {
	is_securitycontextconstraints
	allow_privileged_container := data.sccdenypriv.allowPrivilegedContainer
	input.allowPrivilegedContainer == allow_privileged_container
	msg := sprintf("SCC has allowPrivilegedContainer set to %v", [allow_privileged_container])
}
