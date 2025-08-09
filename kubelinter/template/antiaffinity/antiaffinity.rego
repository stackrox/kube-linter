package kubelinter.template.antiaffinity

import data.kubelinter.objectkinds
import future.keywords.in

deny contains msg if {
	objectkinds.is_deployment_like
	min_replicas := data.antiaffinity.min_replicas
	replicas := replicas_count()
	replicas >= min_replicas
	not has_anti_affinity()
	replicaText := replica_text(replicas)
	msg := sprintf("object has %d %s but does not specify inter pod anti-affinity", [replicas, replicaText])
}

deny contains msg if {
	objectkinds.is_deployment_like
	min_replicas := data.antiaffinity.min_replicas
	replicas := replicas_count()
	replicas >= min_replicas
	has_anti_affinity()
	not has_valid_anti_affinity_rules()
	replicaText := replica_text(replicas)
	msg := sprintf(
		"object has %d %s but does not specify preferred or required inter pod anti-affinity during scheduling",
		[replicas, replicaText],
	)
}

deny contains msg if {
	objectkinds.is_deployment_like
	min_replicas := data.antiaffinity.min_replicas
	replicas := replicas_count()
	replicas >= min_replicas
	has_anti_affinity()
	has_invalid_anti_affinity_config()
	msg := anti_affinity_error_message()
}

replicas_count := input.spec.replicas

replicas_count := 1

replica_text(replicas) := "replicas" if {
	replicas > 1
}

replica_text(1) := "replica"

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
	some term in all_anti_affinity_terms()
	not is_valid_affinity_term(term)
}

all_anti_affinity_terms := terms if {
	preferred := input.spec.template.spec.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution
	required := input.spec.template.spec.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution

	terms := array.concat(
		[term.podAffinityTerm | some term in preferred],
		[term | some term in required],
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

anti_affinity_error_message := msg if {
	some term in all_anti_affinity_terms()
	not is_valid_affinity_term(term)
	msg := sprintf("anti-affinity configuration is invalid for term with topology key %q", [term.topologyKey])
}
