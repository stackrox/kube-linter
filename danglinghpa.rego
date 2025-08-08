package kubelinter.template.danglinghpa

import data.kubelinter.objectkinds.is_horizontalpodautoscaler
import future.keywords.in

deny[msg] {
	is_horizontalpodautoscaler
	target := input.spec.scaleTargetRef
	not target_exists(target)
	msg := sprintf("no resources found matching HorizontalPodAutoscaler scaleTargetRef (%v)", [target])
}

target_exists(target) {
	some deployment in data.objects
	deployment.kind == target.kind
	deployment.metadata.name == target.name
	deployment.apiVersion == target.apiVersion
}