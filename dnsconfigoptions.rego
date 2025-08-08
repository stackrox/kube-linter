package kubelinter.template.dnsconfigoptions

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny[msg] {
	is_deployment_like
	some container in input.spec.template.spec.containers
	some option in container.dnsConfig.options
	not is_allowed_dns_option(option)
	msg := sprintf("container %q has disallowed DNS option %q", [container.name, option.name])
}

deny[msg] {
	is_deployment_like
	some container in input.spec.template.spec.containers
	some option in container.dnsConfig.options
	not has_required_dns_option(option)
	msg := sprintf("container %q is missing required DNS option %q", [container.name, option.name])
}

deny[msg] {
	is_deployment_like
	some container in input.spec.template.spec.containers
	not has_dns_config(container)
	msg := sprintf("container %q has no DNS config specified", [container.name])
}