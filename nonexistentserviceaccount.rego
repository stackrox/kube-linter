package kubelinter.template.nonexistentserviceaccount

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny contains msg if {
	is_deployment_like
	serviceAccountName := input.spec.template.spec.serviceAccountName
	not serviceaccount_exists(serviceAccountName)
	msg := sprintf("service account %q not found", [serviceAccountName])
}

serviceaccount_exists(serviceAccountName) if {
	some sa in data.objects
	sa.kind == "ServiceAccount"
	sa.metadata.name == serviceAccountName
	sa.metadata.namespace == input.metadata.namespace
}