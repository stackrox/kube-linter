package kubelinter.template.memoryrequirements

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	requirementsType := data.memoryrequirements.requirementsType
	lowerBoundMB := data.memoryrequirements.lowerBoundMB
	upperBoundMB := data.memoryrequirements.upperBoundMB

	# Check requests
	is_request_type(requirementsType)
	memoryRequest := container.resources.requests.memory
	memory_bytes := parse_memory_bytes(memoryRequest)
	memory_mb := memory_bytes / (1024 * 1024)
	memory_mb >= lowerBoundMB
	upper_bound_valid(upperBoundMB, memory_mb)
	msg := sprintf("container %q has memory request %s", [container.name, memoryRequest])
}

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	requirementsType := data.memoryrequirements.requirementsType
	lowerBoundMB := data.memoryrequirements.lowerBoundMB
	upperBoundMB := data.memoryrequirements.upperBoundMB

	# Check limits
	is_limit_type(requirementsType)
	memoryLimit := container.resources.limits.memory
	memory_bytes := parse_memory_bytes(memoryLimit)
	memory_mb := memory_bytes / (1024 * 1024)
	memory_mb >= lowerBoundMB
	upper_bound_valid(upperBoundMB, memory_mb)
	msg := sprintf("container %q has memory limit %s", [container.name, memoryLimit])
}

is_request_type(requirementsType) if {
	requirementsType == "request"
}

is_request_type(requirementsType) if {
	requirementsType == "any"
}

is_limit_type(requirementsType) if {
	requirementsType == "limit"
}

is_limit_type(requirementsType) if {
	requirementsType == "any"
}

upper_bound_valid(upperBoundMB, memory_mb) if {
	upperBoundMB == null
}

upper_bound_valid(upperBoundMB, memory_mb) if {
	memory_mb <= upperBoundMB
}

# Helper function to parse memory value to bytes
parse_memory_bytes(memory) := bytes if {
	# Handle "100Mi" format
	regex.match("^([0-9]+)Mi$", memory)
	parts := regex.split("^([0-9]+)Mi$", memory, -1)
	bytes := (to_number(parts[1]) * 1024) * 1024
}

parse_memory_bytes(memory) := bytes if {
	# Handle "100Gi" format
	regex.match("^([0-9]+)Gi$", memory)
	parts := regex.split("^([0-9]+)Gi$", memory, -1)
	bytes := ((to_number(parts[1]) * 1024) * 1024) * 1024
}

parse_memory_bytes(memory) := bytes if {
	# Handle "100Ki" format
	regex.match("^([0-9]+)Ki$", memory)
	parts := regex.split("^([0-9]+)Ki$", memory, -1)
	bytes := to_number(parts[1]) * 1024
}

parse_memory_bytes(memory) := bytes if {
	# Handle "100" format (bytes)
	not regex.match("^([0-9]+)[KMG]i$", memory)
	regex.match("^([0-9]+)$", memory)
	parts := regex.split("^([0-9]+)$", memory, -1)
	bytes := to_number(parts[1])
}
