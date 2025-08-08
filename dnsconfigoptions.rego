package kubelinter.template.dnsconfigoptions

import kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	not input.spec.template.spec.dnsConfig
	msg := "Object does not define any DNSConfig rules."
}

deny contains msg if {
	is_deployment_like
	input.spec.template.spec.dnsConfig
	not input.spec.template.spec.dnsConfig.options
	msg := "Object does not define any DNSConfig Options."
}

deny contains msg if {
	is_deployment_like
	input.spec.template.spec.dnsConfig
	input.spec.template.spec.dnsConfig.options
	key := data.dnsconfigoptions.key
	value := data.dnsconfigoptions.value
	not has_dnsconfig_option(key, value)
	msg := sprintf("DNSConfig Options \"%s:%s\" not found.", [key, value])
}

has_dnsconfig_option(key, value) {
	some option in input.spec.template.spec.dnsConfig.options
	option.name == key
	option.value == value
}