package kubelinter.template.sysctl

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny contains msg if {
	is_deployment_like
	some sysctl in input.spec.template.spec.securityContext.sysctls
	some unsafe_sysctl in data.sysctl.unsafeSysCtls
	startswith(sysctl.name, unsafe_sysctl)
	msg := sprintf("resource specifies unsafe sysctl %q.", [sysctl.name])
}
