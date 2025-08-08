package kubelinter.template.cpurequirements

import kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	requirementsType := data.cpurequirements.requirementsType
	lowerBound := data.cpurequirements.lowerBoundMillis
	upperBound := data.cpurequirements.upperBoundMillis

	# Check requests
	(requirementsType == "request" || requirementsType == "any")
	cpuRequest := container.resources.requests.cpu
	cpu_millis := parse_cpu_millis(cpuRequest)
	cpu_millis >= lowerBound
	(upperBound == null || cpu_millis <= upperBound)
	msg := sprintf("container %q has cpu request %s", [container.name, cpuRequest])
}

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	requirementsType := data.cpurequirements.requirementsType
	lowerBound := data.cpurequirements.lowerBoundMillis
	upperBound := data.cpurequirements.upperBoundMillis

	# Check limits
	(requirementsType == "limit" || requirementsType == "any")
	cpuLimit := container.resources.limits.cpu
	cpu_millis := parse_cpu_millis(cpuLimit)
	cpu_millis >= lowerBound
	(upperBound == null || cpu_millis <= upperBound)
	msg := sprintf("container %q has cpu limit %s", [container.name, cpuLimit])
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