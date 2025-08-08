package kubelinter.template.memoryrequirements

import kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	requirementsType := data.memoryrequirements.requirementsType
	lowerBoundMB := data.memoryrequirements.lowerBoundMB
	upperBoundMB := data.memoryrequirements.upperBoundMB

	# Check requests
	(requirementsType == "request" || requirementsType == "any")
	memoryRequest := container.resources.requests.memory
	memory_bytes := parse_memory_bytes(memoryRequest)
	memory_mb := memory_bytes / (1024 * 1024)
	memory_mb >= lowerBoundMB
	(upperBoundMB == null || memory_mb <= upperBoundMB)
	msg := sprintf("container %q has memory request %s", [container.name, memoryRequest])
}

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	requirementsType := data.memoryrequirements.requirementsType
	lowerBoundMB := data.memoryrequirements.lowerBoundMB
	upperBoundMB := data.memoryrequirements.upperBoundMB

	# Check limits
	(requirementsType == "limit" || requirementsType == "any")
	memoryLimit := container.resources.limits.memory
	memory_bytes := parse_memory_bytes(memoryLimit)
	memory_mb := memory_bytes / (1024 * 1024)
	memory_mb >= lowerBoundMB
	(upperBoundMB == null || memory_mb <= upperBoundMB)
	msg := sprintf("container %q has memory limit %s", [container.name, memoryLimit])
}

# Helper function to parse memory value to bytes
parse_memory_bytes(memory) := bytes {
	# Handle "100Mi" format
	regex.match("^([0-9]+)Mi$", memory)
	parts := regex.split("^([0-9]+)Mi$", memory, -1)
	bytes := to_number(parts[1]) * 1024 * 1024
}

parse_memory_bytes(memory) := bytes {
	# Handle "100Gi" format
	regex.match("^([0-9]+)Gi$", memory)
	parts := regex.split("^([0-9]+)Gi$", memory, -1)
	bytes := to_number(parts[1]) * 1024 * 1024 * 1024
}

parse_memory_bytes(memory) := bytes {
	# Handle "100Ki" format
	regex.match("^([0-9]+)Ki$", memory)
	parts := regex.split("^([0-9]+)Ki$", memory, -1)
	bytes := to_number(parts[1]) * 1024
}

parse_memory_bytes(memory) := bytes {
	# Handle "100" format (bytes)
	not regex.match("^([0-9]+)[KMG]i$", memory)
	regex.match("^([0-9]+)$", memory)
	parts := regex.split("^([0-9]+)$", memory, -1)
	bytes := to_number(parts[1])
}