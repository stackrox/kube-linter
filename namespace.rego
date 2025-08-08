package kubelinter.template.namespace

import kubelinter.objectkinds.is_deployment_like
import kubelinter.objectkinds.is_service

deny contains msg if {
	(is_deployment_like || is_service)
	namespace := input.metadata.namespace
	namespace == "default"
	msg := "object in default namespace"
}

deny contains msg if {
	(is_deployment_like || is_service)
	not input.metadata.namespace
	msg := "object in default namespace"
}