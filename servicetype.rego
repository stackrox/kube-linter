package kubelinter.template.servicetype

import kubelinter.objectkinds.is_service

deny contains msg if {
	is_service
	some forbiddenType in data.servicetype.forbiddenServiceTypes
	input.spec.type == forbiddenType
	msg := sprintf("%q service type is forbidden.", [forbiddenType])
}