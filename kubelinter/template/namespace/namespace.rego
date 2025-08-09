package kubelinter.template.namespace

import data.kubelinter.objectkinds.is_deployment_like
import data.kubelinter.objectkinds.is_service
import future.keywords.in

deny contains msg if {
	is_deployment_like_or_service()
	namespace := input.metadata.namespace
	namespace == "default"
	msg := "object in default namespace"
}

deny contains msg if {
	is_deployment_like_or_service()
	not input.metadata.namespace
	msg := "object in default namespace"
}

is_deployment_like_or_service if {
	is_deployment_like
}

is_deployment_like_or_service if {
	is_service
}
