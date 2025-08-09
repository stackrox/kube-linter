package abc.xyz

deny contains msg if {
	value := input.request.object.metadata.name

	# Check if the label value is formatted correctly.
	not startswith(value, "asdfasdf-")

	# Construct an error message to return to the user.
	msg := sprintf("3; found %v", [value])
}
