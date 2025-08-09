package kubelinter.template.nonexistentserviceaccount

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny contains msg if {
	is_deployment_like
	service_account_name := input.spec.template.spec.serviceAccountName
	not serviceaccount_exists(service_account_name)
	msg := sprintf("service account %q not found", [service_account_name])
}

serviceaccount_exists(service_account_name) if {
	some sa in data.objects
	sa.kind == "ServiceAccount"
	sa.metadata.name == service_account_name
	sa.metadata.namespace == input.metadata.namespace
}
