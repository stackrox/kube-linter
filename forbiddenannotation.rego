package kubelinter.template.forbiddenannotation

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny[msg] {
	is_deployment_like
	some annotation in input.metadata.annotations
	not is_allowed_annotation(annotation)
	msg := sprintf("annotation %q is forbidden", [annotation])
}