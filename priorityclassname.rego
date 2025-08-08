package kubelinter.template.priorityclassname

import kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	priorityClassName := input.spec.template.spec.priorityClassName
	priorityClassName != ""
	not is_accepted_priority_class(priorityClassName)
	acceptedClasses := data.priorityclassname.acceptedPriorityClassNames
	msg := sprintf("object has a priority class name defined with '%s' but the only accepted priority class names are '%s'", [priorityClassName, array.join(acceptedClasses, ", ")])
}

is_accepted_priority_class(priorityClassName) {
	some acceptedClass in data.priorityclassname.acceptedPriorityClassNames
	priorityClassName == acceptedClass
}