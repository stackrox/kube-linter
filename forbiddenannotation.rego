package kubelinter.template.forbiddenannotation

import kubelinter.objectkinds.is_any

deny contains msg if {
	is_any
	key := data.forbiddenannotation.key
	value := data.forbiddenannotation.value
	has_annotation(key, value)
	msg := sprintf("object has forbidden annotation %q with value %q", [key, value])
}

has_annotation(key, value) {
	input.metadata.annotations[key] == value
}