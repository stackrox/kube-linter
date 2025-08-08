package kubelinter.template.requiredlabel

import kubelinter.objectkinds.is_any

deny contains msg if {
	is_any
	key := data.requiredlabel.key
	value := data.requiredlabel.value
	not has_label(key, value)
	msg := sprintf("object is missing required label %q with value %q", [key, value])
}

has_label(key, value) {
	input.metadata.labels[key] == value
}