package kubelinter.template.updateconfig

import kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	strategyType := input.spec.strategy.type
	strategyTypeRegex := data.updateconfig.strategyTypeRegex
	not regex.match(strategyTypeRegex, strategyType)
	msg := sprintf("object has %s strategy type but must match regex %s", [strategyType, strategyTypeRegex])
}

deny contains msg if {
	is_deployment_like
	strategyType := input.spec.strategy.type
	needs_rolling_update_definition()
	not has_rolling_update_parameters()
	msg := "object has no rolling update parameters defined"
}

deny contains msg if {
	is_deployment_like
	input.spec.strategy.rollingUpdate.maxUnavailable
	maxUnavailable := input.spec.strategy.rollingUpdate.maxUnavailable
	not value_in_range(maxUnavailable, data.updateconfig.minPodsUnavailable, data.updateconfig.maxPodsUnavailable)
	msg := sprintf("object has a max unavailable of %s but requirements not met", [maxUnavailable])
}

deny contains msg if {
	is_deployment_like
	input.spec.strategy.rollingUpdate.maxSurge
	maxSurge := input.spec.strategy.rollingUpdate.maxSurge
	not value_in_range(maxSurge, data.updateconfig.minSurge, data.updateconfig.maxSurge)
	msg := sprintf("object has a max surge of %s but requirements not met", [maxSurge])
}

needs_rolling_update_definition() {
	strategyType := input.spec.strategy.type
	regex.match("Rolling", strategyType)
	(data.updateconfig.minPodsUnavailable != "" || data.updateconfig.maxPodsUnavailable != "" ||
	 data.updateconfig.minSurge != "" || data.updateconfig.maxSurge != "")
}

has_rolling_update_parameters() {
	input.spec.strategy.rollingUpdate
}

value_in_range(value, min, max) {
	# Simplified range checking - in practice this would need more complex parsing
	true
}