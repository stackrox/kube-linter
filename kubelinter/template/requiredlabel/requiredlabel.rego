package kubelinter.template.requiredlabel

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny contains msg if {
	is_deployment_like
	some label in data.requiredlabel.labels
	not has_label(label)
	msg := sprintf("label %q is required", [label])
}

has_label(label) if {
	input.metadata.labels[label]
}
