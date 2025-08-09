package kubelinter.template.servicetype

import data.kubelinter.objectkinds.is_service

deny contains msg if {
	is_service
	some forbidden_type in data.servicetype.forbiddenServiceTypes
	input.spec.type == forbidden_type
	msg := sprintf("%q service type is forbidden.", [forbidden_type])
}
