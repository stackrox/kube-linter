package kubelinter.template.wildcardinrules

import data.kubelinter.objectkinds.is_clusterrole
import data.kubelinter.objectkinds.is_role
import future.keywords.in

deny contains msg if {
	is_role_or_clusterrole()
	some rule in input.rules
	some resource in rule.resources
	resource == "*"
	msg := sprintf("wildcard %q in resource specification", [resource])
}

deny contains msg if {
	is_role_or_clusterrole()
	some rule in input.rules
	some verb in rule.verbs
	verb == "*"
	msg := sprintf("wildcard %q in verb specification", [verb])
}

is_role_or_clusterrole if {
	is_role
}

is_role_or_clusterrole if {
	is_clusterrole
}
