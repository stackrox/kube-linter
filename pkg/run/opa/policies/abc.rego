package abc

deny contains msg if {
	value := input.request.object.metadata.labels.costcenter

	# Check if the label value is formatted correctly.
	not startswith(value, "cccode-")

	# Construct an error message to return to the user.
	msg := sprintf("1 Costcenter code must start with cccode-; found %v", [value])
}

deny contains msg if {
	value := input.request.object.metadata.name

	# Check if the label value is formatted correctly.
	not startswith(value, data.name)

	# Construct an error message to return to the user.
	msg := sprintf("2 name must start with cccode-; found %v", [data.name])
}


deny contains msg if {
	not is_ingress
	msg := "not ingres"
}
