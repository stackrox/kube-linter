package kubelinter.template.requiredannotation

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny contains msg if {
	is_deployment_like
	some annotation in data.requiredannotation.annotations
	not has_annotation(annotation)
	msg := sprintf("annotation %q is required", [annotation])
}

has_annotation(annotation) if {
	input.metadata.annotations[annotation]
}