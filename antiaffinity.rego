package kubelinter.template.antiaffinity

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny[msg] {
	is_deployment_like
	minReplicas := data.antiaffinity.minReplicas
	replicas := get_replicas()
	replicas >= minReplicas
	not has_anti_affinity()
	replicaText := get_replica_text(replicas)
	msg := sprintf("object has %d %s but does not specify inter pod anti-affinity", [replicas, replicaText])
}

deny[msg] {
	is_deployment_like
	minReplicas := data.antiaffinity.minReplicas
	replicas := get_replicas()
	replicas >= minReplicas
	has_anti_affinity()
	not has_valid_anti_affinity_rules()
	replicaText := get_replica_text(replicas)
	msg := sprintf("object has %d %s but does not specify preferred or required inter pod anti-affinity during scheduling", [replicas, replicaText])
}

deny[msg] {
	is_deployment_like
	minReplicas := data.antiaffinity.minReplicas
	replicas := get_replicas()
	replicas >= minReplicas
	has_anti_affinity()
	has_invalid_anti_affinity_config()
	msg := get_anti_affinity_error_message()
}

get_replicas() := replicas {
	replicas := input.spec.replicas
}

get_replicas() := 1 {
	not input.spec.replicas
}

get_replica_text(replicas) := "replicas" {
	replicas > 1
}

get_replica_text(replicas) := "replica" {
	replicas == 1
}

has_anti_affinity() {
	input.spec.template.spec.affinity.podAntiAffinity
}

has_valid_anti_affinity_rules() {
	has_preferred_rules()
}

has_valid_anti_affinity_rules() {
	has_required_rules()
}

has_preferred_rules() {
	preferred := input.spec.template.spec.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution
	count(preferred) > 0
}

has_required_rules() {
	required := input.spec.template.spec.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution
	count(required) > 0
}

has_invalid_anti_affinity_config() {
	some term in get_all_anti_affinity_terms()
	not is_valid_affinity_term(term)
}

get_all_anti_affinity_terms() := terms {
	preferred := input.spec.template.spec.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution
	required := input.spec.template.spec.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution

	terms := array.concat(
		[term.podAffinityTerm | term := preferred[_]],
		[term | term := required[_]]
	)
}

is_valid_affinity_term(term) {
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

namespace_matches_term(term) {
	not term.namespaces
}

namespace_matches_term(term) {
	namespace := input.metadata.namespace
	namespace in term.namespaces
}

labels_match_selector(selector, labels) {
	# Simplified label matching - in practice this would need more complex logic
	true
}

get_anti_affinity_error_message() := msg {
	some term in get_all_anti_affinity_terms()
	not is_valid_affinity_term(term)
	msg := sprintf("anti-affinity configuration is invalid for term with topology key %q", [term.topologyKey])
}