package kubelinter.template.cpurequirements

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny[msg] {
	is_deployment_like
	some container in input.spec.template.spec.containers
	requirementsType := data.cpurequirements.requirementsType
	lowerBound := data.cpurequirements.lowerBoundMillis
	upperBound := data.cpurequirements.upperBoundMillis

	# Check requests
	is_request_type(requirementsType)
	cpuRequest := container.resources.requests.cpu
	cpu_millis := parse_cpu_millis(cpuRequest)
	cpu_millis >= lowerBound
	upper_bound_valid(upperBound, cpu_millis)
	msg := sprintf("container %q has cpu request %s", [container.name, cpuRequest])
}

deny[msg] {
	is_deployment_like
	some container in input.spec.template.spec.containers
	requirementsType := data.cpurequirements.requirementsType
	lowerBound := data.cpurequirements.lowerBoundMillis
	upperBound := data.cpurequirements.upperBoundMillis

	# Check limits
	is_limit_type(requirementsType)
	cpuLimit := container.resources.limits.cpu
	cpu_millis := parse_cpu_millis(cpuLimit)
	cpu_millis >= lowerBound
	upper_bound_valid(upperBound, cpu_millis)
	msg := sprintf("container %q has cpu limit %s", [container.name, cpuLimit])
}

is_request_type(requirementsType) {
	requirementsType == "request"
}

is_request_type(requirementsType) {
	requirementsType == "any"
}

is_limit_type(requirementsType) {
	requirementsType == "limit"
}

is_limit_type(requirementsType) {
	requirementsType == "any"
}

upper_bound_valid(upperBound, cpu_millis) {
	upperBound == null
}

upper_bound_valid(upperBound, cpu_millis) {
	cpu_millis <= upperBound
}

# Helper function to parse CPU value to millicores
parse_cpu_millis(cpu) := millis {
	# Handle "100m" format (millicores)
	regex.match("^([0-9]+)m$", cpu)
	parts := regex.split("^([0-9]+)m$", cpu, -1)
	millis := to_number(parts[1])
}

parse_cpu_millis(cpu) := millis {
	# Handle "1" format (cores) - convert to millicores
	not regex.match("^([0-9]+)m$", cpu)
	not regex.match("^([0-9]+\\.[0-9]+)$", cpu)
	regex.match("^([0-9]+)$", cpu)
	parts := regex.split("^([0-9]+)$", cpu, -1)
	millis := to_number(parts[1]) * 1000
}

parse_cpu_millis(cpu) := millis {
	# Handle "1.5" format (fractional cores) - convert to millicores
	not regex.match("^([0-9]+)m$", cpu)
	regex.match("^([0-9]+\\.[0-9]+)$", cpu)
	parts := regex.split("^([0-9]+)\\.([0-9]+)$", cpu, -1)
	whole := to_number(parts[1])
	fraction := to_number(parts[2])
	millis := (whole * 1000) + (fraction * 1000)
}