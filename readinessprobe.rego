package kubelinter.template.readinessprobe

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny contains msg if {
	is_deployment_like
	some container in input.spec.template.spec.containers
	not container.readinessProbe
	msg := sprintf("container %q does not specify a readiness probe", [container.name])
}
