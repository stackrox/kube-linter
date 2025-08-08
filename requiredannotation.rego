package kubelinter.template.requiredannotation

import kubelinter.objectkinds.is_any

deny contains msg if {
	is_any
	key := data.requiredannotation.key
	value := data.requiredannotation.value
	not has_annotation(key, value)
	msg := sprintf("object is missing required annotation %q with value %q", [key, value])
}

has_annotation(key, value) {
	input.metadata.annotations[key] == value
}