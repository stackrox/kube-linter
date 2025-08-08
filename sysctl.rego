package kubelinter.template.sysctl

import kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	some sysctl in input.spec.template.spec.securityContext.sysctls
	some unsafeSysctl in data.sysctl.unsafeSysCtls
	startswith(sysctl.name, unsafeSysctl)
	msg := sprintf("resource specifies unsafe sysctl %q.", [sysctl.name])
}