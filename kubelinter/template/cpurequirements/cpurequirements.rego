package kubelinter.template.cpurequirements

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	requirements_type := data.cpurequirements.requirementsType

	# Check requests
	is_request_type(requirements_type)
	cpu_request := container.resources.requests.cpu
	cpu_millis := parse_cpu_millis(cpu_request)
	cpu_millis >= data.cpurequirements.lowerBoundMillis
	upper_bound_valid(data.cpurequirements.upperBoundMillis, cpu_millis)
	msg := sprintf("container %q has cpu request %s", [container.name, cpu_request])
}

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	requirements_type := data.cpurequirements.requirementsType

	# Check limits
	is_limit_type(requirements_type)
	cpu_limit := container.resources.limits.cpu
	cpu_millis := parse_cpu_millis(cpu_limit)
	cpu_millis >= data.cpurequirements.lowerBoundMillis
	upper_bound_valid(data.cpurequirements.upperBoundMillis, cpu_millis)
	msg := sprintf("container %q has cpu limit %s", [container.name, cpu_limit])
}

is_request_type("request")

is_request_type("any")

is_limit_type("limit")

is_limit_type("any")

upper_bound_valid(null, cpu_millis)

upper_bound_valid(upper_bound, cpu_millis) if {
	cpu_millis <= upper_bound
}

# Helper function to parse CPU value to millicores
parse_cpu_millis(cpu) := to_number(regex.find_n(`[0-9]+`, cpu, 1)[0]) if {
	# Handle "100m" format (millicores)
	regex.match(`^([0-9]+)m$`, cpu)
}

parse_cpu_millis(cpu) := to_number(regex.find_n(`[0-9]+`, cpu, 1)[0]) * 1000 if {
	# Handle "1" format (cores) - convert to millicores
	not regex.match(`^([0-9]+)m$`, cpu)
	not regex.match(`^([0-9]+\.[0-9]+)$`, cpu)
	regex.match(`^([0-9]+)$`, cpu)
}

parse_cpu_millis(cpu) := to_number(cpu) * 1000 if {
	# Handle "1.5" format (fractional cores) - convert to millicores
	not regex.match(`^([0-9]+)m$`, cpu)
	regex.match(`^([0-9]+\.[0-9]+)$`, cpu)
}
