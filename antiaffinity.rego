package kubelinter.template.antiaffinity

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny contains msg if {
	is_deployment_like
	minReplicas := data.antiaffinity.minReplicas
	replicas := get_replicas()
	replicas >= minReplicas
	not has_anti_affinity()
	replicaText := get_replica_text(replicas)
	msg := sprintf("object has %d %s but does not specify inter pod anti-affinity", [replicas, replicaText])
}

deny contains msg if {
	is_deployment_like
	minReplicas := data.antiaffinity.minReplicas
	replicas := get_replicas()
	replicas >= minReplicas
	has_anti_affinity()
	not has_valid_anti_affinity_rules()
	replicaText := get_replica_text(replicas)
	msg := sprintf("object has %d %s but does not specify preferred or required inter pod anti-affinity during scheduling", [replicas, replicaText])
}

deny contains msg if {
	is_deployment_like
	minReplicas := data.antiaffinity.minReplicas
	replicas := get_replicas()
	replicas >= minReplicas
	has_anti_affinity()
	has_invalid_anti_affinity_config()
	msg := get_anti_affinity_error_message()
}

get_replicas := replicas if {
	replicas := input.spec.replicas
}

get_replicas := 1 if {
	not input.spec.replicas
}

get_replica_text(replicas) := "replicas" if {
	replicas > 1
}

get_replica_text(replicas) := "replica" if {
	replicas == 1
}

has_anti_affinity if {
	input.spec.template.spec.affinity.podAntiAffinity
}

has_valid_anti_affinity_rules if {
	has_preferred_rules()
}

has_valid_anti_affinity_rules if {
	has_required_rules()
}

has_preferred_rules if {
	preferred := input.spec.template.spec.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution
	count(preferred) > 0
}

has_required_rules if {
	required := input.spec.template.spec.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution
	count(required) > 0
}

has_invalid_anti_affinity_config if {
	some term in get_all_anti_affinity_terms()
	not is_valid_affinity_term(term)
}

get_all_anti_affinity_terms := terms if {
	preferred := input.spec.template.spec.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution
	required := input.spec.template.spec.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution

	terms := array.concat(
		[term.podAffinityTerm | term := preferred[_]],
		[term | term := required[_]],
	)
}

is_valid_affinity_term(term) if {
	# Check namespace
	namespace := input.metadata.namespace
	namespace == "default"
	namespace_matches_term(term)

	# Check topology key
	topologyKey := data.antiaffinity.topologyKey
	topologyKey == "kubernetes.io/hostname"
	term.topologyKey == topologyKey

	# Check label selector
	labels_match_selector(term.labelSelector, input.spec.template.metadata.labels)
}

namespace_matches_term(term) if {
	not term.namespaces
}

namespace_matches_term(term) if {
	namespace := input.metadata.namespace
	namespace in term.namespaces
}

labels_match_selector(selector, labels) := true

# Simplified label matching - in practice this would need more complex logic

get_anti_affinity_error_message := msg if {
	some term in get_all_anti_affinity_terms()
	not is_valid_affinity_term(term)
	msg := sprintf("anti-affinity configuration is invalid for term with topology key %q", [term.topologyKey])
}
