package kubelinter.template.sccdenypriv

import kubelinter.objectkinds.is_securitycontextconstraints

deny contains msg if {
	is_securitycontextconstraints
	allowPrivilegedContainer := data.sccdenypriv.allowPrivilegedContainer
	input.allowPrivilegedContainer == allowPrivilegedContainer
	msg := sprintf("SCC has allowPrivilegedContainer set to %v", [allowPrivilegedContainer])
}