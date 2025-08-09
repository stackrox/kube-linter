package kubelinter.template.memoryrequirements

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	requirements_type := data.memoryrequirements.requirementsType

	# Check requests
	is_request_type(requirements_type)
	memory_request := container.resources.requests.memory
	memory_bytes := parse_memory_bytes(memory_request)
	memory_mb := memory_bytes / (1024 * 1024)
	memory_mb >= data.memoryrequirements.lowerBoundMB
	upper_bound_valid(data.memoryrequirements.upperBoundMB, memory_mb)
	msg := sprintf("container %q has memory request %s", [container.name, memory_request])
}

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	requirements_type := data.memoryrequirements.requirementsType

	# Check limits
	is_limit_type(requirements_type)
	memory_limit := container.resources.limits.memory
	memory_bytes := parse_memory_bytes(memory_limit)
	memory_mb := memory_bytes / (1024 * 1024)
	memory_mb >= data.memoryrequirements.lowerBoundMB
	upper_bound_valid(data.memoryrequirements.upperBoundMB, memory_mb)
	msg := sprintf("container %q has memory limit %s", [container.name, memory_limit])
}

is_request_type("request")

is_request_type("any")

is_limit_type("limit")

is_limit_type("any")

upper_bound_valid(null, memory_mb)

upper_bound_valid(upper_bound_mb, memory_mb) if {
	memory_mb <= upper_bound_mb
}

# Helper function to parse memory value to bytes
parse_memory_bytes(memory) := (to_number(regex.find_n(`[0-9]+`, memory, 1)[0]) * 1024) * 1024 if {
	# Handle "100Mi" format
	regex.match(`^([0-9]+)Mi$`, memory)
}

parse_memory_bytes(memory) := ((to_number(regex.find_n(`[0-9]+`, memory, 1)[0]) * 1024) * 1024) * 1024 if {
	# Handle "100Gi" format
	regex.match(`^([0-9]+)Gi$`, memory)
}

parse_memory_bytes(memory) := to_number(regex.find_n(`[0-9]+`, memory, 1)[0]) * 1024 if {
	# Handle "100Ki" format
	regex.match(`^([0-9]+)Ki$`, memory)
}

parse_memory_bytes(memory) := to_number(memory) if {
	# Handle "100" format (bytes)
	not regex.match(`^([0-9]+)[KMG]i$`, memory)
	regex.match(`^([0-9]+)$`, memory)
}
