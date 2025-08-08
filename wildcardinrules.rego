package kubelinter.template.wildcardinrules

import kubelinter.objectkinds.is_role
import kubelinter.objectkinds.is_clusterrole

deny contains msg if {
	(is_role || is_clusterrole)
	some rule in input.rules
	some resource in rule.resources
	resource == "*"
	msg := sprintf("wildcard %q in resource specification", [resource])
}

deny contains msg if {
	(is_role || is_clusterrole)
	some rule in input.rules
	some verb in rule.verbs
	verb == "*"
	msg := sprintf("wildcard %q in verb specification", [verb])
}