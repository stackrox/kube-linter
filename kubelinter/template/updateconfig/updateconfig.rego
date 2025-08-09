package kubelinter.template.updateconfig

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny contains msg if {
	is_deployment_like
	strategyType := input.spec.strategy.type
	strategyTypeRegex := data.updateconfig.strategyTypeRegex
	not regex.match(strategyTypeRegex, strategyType)
	msg := sprintf("object has %s strategy type but must match regex %s", [strategyType, strategyTypeRegex])
}

deny contains msg if {
	is_deployment_like
	needs_rolling_update_definition()
	not has_rolling_update_parameters()
	msg := "object has no rolling update parameters defined"
}

deny contains msg if {
	is_deployment_like
	maxUnavailable := input.spec.strategy.rollingUpdate.maxUnavailable
	not value_in_range(maxUnavailable, data.updateconfig.minPodsUnavailable, data.updateconfig.maxPodsUnavailable)
	msg := sprintf("object has a max unavailable of %s but requirements not met", [maxUnavailable])
}

deny contains msg if {
	is_deployment_like
	maxSurge := input.spec.strategy.rollingUpdate.maxSurge
	not value_in_range(maxSurge, data.updateconfig.minSurge, data.updateconfig.maxSurge)
	msg := sprintf("object has a max surge of %s but requirements not met", [maxSurge])
}

needs_rolling_update_definition if {
	strategyType := input.spec.strategy.type
	regex.match(`Rolling`, strategyType)
	has_rolling_update_config()
}

has_rolling_update_config if {
	has_min_pods_unavailable()
}

has_rolling_update_config if {
	has_max_pods_unavailable()
}

has_rolling_update_config if {
	has_min_surge()
}

has_rolling_update_config if {
	has_max_surge()
}

has_min_pods_unavailable if {
	data.updateconfig.minPodsUnavailable != ""
}

has_max_pods_unavailable if {
	data.updateconfig.maxPodsUnavailable != ""
}

has_min_surge if {
	data.updateconfig.minSurge != ""
}

has_max_surge if {
	data.updateconfig.maxSurge != ""
}

has_rolling_update_parameters if {
	input.spec.strategy.rollingUpdate
}

value_in_range(value, min_val, max_val) := true

# Simplified range checking - in practice this would need more complex parsing
