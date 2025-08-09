package kubelinter.template.accesstoresources

import data.kubelinter.objectkinds

import future.keywords.every
import future.keywords.in

deny contains msg if {
	objectkinds.is_rolebinding
	role_name := input.roleRef.name
	namespace := input.metadata.namespace
	namespace := "default"
	not role_exists(role_name, namespace)
	data.accesstoresources.flagRolesNotFound == true
	msg := sprintf("role %q in namespace %q not found", [role_name, namespace])
}

deny contains msg if {
	objectkinds.is_rolebinding
	role_name := input.roleRef.name
	namespace := input.metadata.namespace
	namespace := "default"
	role := find_role(role_name, namespace)
	accesses := check_role_access(role)
	count(accesses) > 0
	msg := sprintf("binding to %q role that has %s", [role_name, array.join(accesses, ", ")])
}

deny contains msg if {
	objectkinds.is_clusterrolebinding
	cluster_role_name := input.roleRef.name
	not clusterrole_exists(cluster_role_name)
	data.accesstoresources.flagRolesNotFound == true
	msg := sprintf("clusterrole %q not found", [cluster_role_name])
}

deny contains msg if {
	objectkinds.is_clusterrolebinding
	cluster_role_name := input.roleRef.name
	cluster_role := find_clusterrole(cluster_role_name)
	accesses := check_clusterrole_access(cluster_role)
	count(accesses) > 0
	msg := sprintf("binding to %q clusterrole that has %s", [cluster_role_name, array.join(accesses, ", ")])
}

deny contains msg if {
	objectkinds.is_clusterrolebinding
	cluster_role_name := input.roleRef.name
	cluster_role := find_clusterrole(cluster_role_name)
	cluster_role.aggregationRule
	some selector in cluster_role.aggregationRule.clusterRoleSelectors
	aggregated_role := find_aggregated_clusterrole(selector)
	accesses := check_clusterrole_access(aggregated_role)
	count(accesses) > 0
	msg := sprintf(
		"binding via aggregationRule to %q clusterrole that has %s",
		[aggregated_role.metadata.name, array.join(accesses, ", ")],
	)
}

# Helper functions
role_exists(role_name, namespace) if {
	some role in data.objects
	role.kind == "Role"
	role.metadata.name == role_name
	role.metadata.namespace == namespace
}

clusterrole_exists(cluster_role_name) if {
	some cluster_role in data.objects
	cluster_role.kind == "ClusterRole"
	cluster_role.metadata.name == cluster_role_name
}

find_role(role_name, namespace) := role if {
	some role in data.objects
	role.kind == "Role"
	role.metadata.name == role_name
	role.metadata.namespace == namespace
}

find_clusterrole(cluster_role_name) := cluster_role if {
	some cluster_role in data.objects
	cluster_role.kind == "ClusterRole"
	cluster_role.metadata.name == cluster_role_name
}

find_aggregated_clusterrole(selector) := cluster_role if {
	some cluster_role in data.objects
	cluster_role.kind == "ClusterRole"
	labels_match(selector.matchLabels, cluster_role.metadata.labels)
}

labels_match(selector_labels, object_labels) if {
	every key, value in selector_labels {
		object_labels[key] == value
	}
}

check_role_access(role) := accesses if {
	some rule in role.rules
	some resource in rule.resources
	some verb in rule.verbs
	resource_matches(resource)
	verb_matches(verb)
	accesses := [sprintf("%v access to %v", [verb, resource])]
}

check_clusterrole_access(cluster_role) := accesses if {
	some rule in cluster_role.rules
	some resource in rule.resources
	some verb in rule.verbs
	resource_matches(resource)
	verb_matches(verb)
	accesses := [sprintf("%v access to %v", [verb, resource])]
}

resource_matches(resource) if {
	some pattern in data.accesstoresources.resources
	pattern == "*"
	resource == "*"
}

resource_matches(resource) if {
	some pattern in data.accesstoresources.resources
	regex.match(pattern, resource)
}

verb_matches(verb) if {
	some pattern in data.accesstoresources.verbs
	pattern == "*"
	verb == "*"
}

verb_matches(verb) if {
	some pattern in data.accesstoresources.verbs
	regex.match(pattern, verb)
}

resource_is_wildcard("*")

verb_is_wildcard("*")
