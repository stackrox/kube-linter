package kubelinter.template.danglingservicemonitor

import kubelinter.objectkinds.is_servicemonitor

deny contains msg if {
	is_servicemonitor
	not has_selector()
	msg := "service monitor has no selector specified"
}

deny contains msg if {
	is_servicemonitor
	has_selector()
	not has_matching_service()
	msg := sprintf("no services found matching the service monitor's label selector (%s) and namespace selector (%s)", [input.spec.selector, input.spec.namespaceSelector])
}

has_selector() {
	nsSelector := input.spec.namespaceSelector
	labelSelectors := input.spec.selector.matchLabels
	(count(nsSelector.matchNames) > 0 || nsSelector.any) || count(labelSelectors) > 0
}

has_matching_service() {
	some service in data.objects
	service.kind == "Service"
	namespace_matches(service)
	labels_match(service)
}

namespace_matches(service) {
	nsSelector := input.spec.namespaceSelector
	nsSelector.any
}

namespace_matches(service) {
	nsSelector := input.spec.namespaceSelector
	service.metadata.namespace in nsSelector.matchNames
}

labels_match(service) {
	labelSelectors := input.spec.selector.matchLabels
	count(labelSelectors) == 0
}

labels_match(service) {
	labelSelectors := input.spec.selector.matchLabels
	count(labelSelectors) > 0
	every key, value in labelSelectors {
		service.metadata.labels[key] == value
	}
}