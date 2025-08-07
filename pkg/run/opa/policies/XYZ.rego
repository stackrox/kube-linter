package abc

deny contains msg if {
	value := input.request.object.metadata.name

	# Check if the label value is formatted correctly.
	not startswith(value, "asdfasdf-")

	# Construct an error message to return to the user.
	msg := sprintf("3; found %v", [value])
}


is_ingress if {
	input.request.kind.kind == "Ingress"
	input.request.kind.group == "extensions"
	input.request.kind.version == "v1beta1"
}