package kubelinter.template.servicetype

import data.kubelinter.objectkinds.is_service

deny contains msg if {
	is_service
	some forbiddenType in data.servicetype.forbiddenServiceTypes
	input.spec.type == forbiddenType
	msg := sprintf("%q service type is forbidden.", [forbiddenType])
}