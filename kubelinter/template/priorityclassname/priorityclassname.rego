package kubelinter.template.priorityclassname

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny contains msg if {
	is_deployment_like
	not input.spec.template.spec.priorityClassName
	msg := "object has no priority class name specified"
}
