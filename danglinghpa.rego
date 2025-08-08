package kubelinter.template.danglinghpa

import kubelinter.objectkinds.is_horizontalpodautoscaler

deny contains msg if {
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