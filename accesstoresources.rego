package kubelinter.template.accesstoresources

import kubelinter.objectkinds.is_role
import kubelinter.objectkinds.is_clusterrole
import kubelinter.objectkinds.is_rolebinding
import kubelinter.objectkinds.is_clusterrolebinding

deny contains msg if {
	is_rolebinding
	roleName := input.roleRef.name
	namespace := input.metadata.namespace
	namespace := "default"
	flagRolesNotFound := data.accesstoresources.flagRolesNotFound
	not role_exists(roleName, namespace)
	flagRolesNotFound == true
	msg := sprintf("role %q in namespace %q not found", [roleName, namespace])
}

deny contains msg if {
	is_rolebinding
	roleName := input.roleRef.name
	namespace := input.metadata.namespace
	namespace := "default"
	role := find_role(roleName, namespace)
	accesses := check_role_access(role)
	count(accesses) > 0
	msg := sprintf("binding to %q role that has %s", [roleName, array.join(accesses, ", ")])
}

deny contains msg if {
	is_clusterrolebinding
	clusterRoleName := input.roleRef.name
	flagRolesNotFound := data.accesstoresources.flagRolesNotFound
	not clusterrole_exists(clusterRoleName)
	flagRolesNotFound == true
	msg := sprintf("clusterrole %q not found", [clusterRoleName])
}

deny contains msg if {
	is_clusterrolebinding
	clusterRoleName := input.roleRef.name
	clusterRole := find_clusterrole(clusterRoleName)
	accesses := check_clusterrole_access(clusterRole)
	count(accesses) > 0
	msg := sprintf("binding to %q clusterrole that has %s", [clusterRoleName, array.join(accesses, ", ")])
}

deny contains msg if {
	is_clusterrolebinding
	clusterRoleName := input.roleRef.name
	clusterRole := find_clusterrole(clusterRoleName)
	clusterRole.aggregationRule
	some selector in clusterRole.aggregationRule.clusterRoleSelectors
	aggregatedRole := find_aggregated_clusterrole(selector)
	accesses := check_clusterrole_access(aggregatedRole)
	count(accesses) > 0
	msg := sprintf("binding via aggregationRule to %q clusterrole that has %s", [aggregatedRole.metadata.name, array.join(accesses, ", ")])
}

# Helper functions
role_exists(roleName, namespace) {
	some role in data.objects
	role.kind == "Role"
	role.metadata.name == roleName
	role.metadata.namespace == namespace
}

clusterrole_exists(clusterRoleName) {
	some clusterRole in data.objects
	clusterRole.kind == "ClusterRole"
	clusterRole.metadata.name == clusterRoleName
}

find_role(roleName, namespace) := role {
	some role in data.objects
	role.kind == "Role"
	role.metadata.name == roleName
	role.metadata.namespace == namespace
}

find_clusterrole(clusterRoleName) := clusterRole {
	some clusterRole in data.objects
	clusterRole.kind == "ClusterRole"
	clusterRole.metadata.name == clusterRoleName
}

find_aggregated_clusterrole(selector) := clusterRole {
	some clusterRole in data.objects
	clusterRole.kind == "ClusterRole"
	labels_match(selector.matchLabels, clusterRole.metadata.labels)
}

labels_match(selectorLabels, objectLabels) {
	every key, value in selectorLabels {
		objectLabels[key] == value
	}
}

check_role_access(role) := accesses {
	some rule in role.rules
	some resource in rule.resources
	some verb in rule.verbs
	resource_matches(resource)
	verb_matches(verb)
	accesses := [sprintf("%v access to %v", [verb, resource])]
}

check_clusterrole_access(clusterRole) := accesses {
	some rule in clusterRole.rules
	some resource in rule.resources
	some verb in rule.verbs
	resource_matches(resource)
	verb_matches(verb)
	accesses := [sprintf("%v access to %v", [verb, resource])]
}

resource_matches(resource) {
	some pattern in data.accesstoresources.resources
	resource == "*" || regex.match(pattern, resource)
}

verb_matches(verb) {
	some pattern in data.accesstoresources.verbs
	verb == "*" || regex.match(pattern, verb)
}