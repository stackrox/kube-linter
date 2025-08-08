package kubelinter.template.nonexistentserviceaccount

import kubelinter.objectkinds.is_deployment_like

deny contains msg if {
	is_deployment_like
	serviceAccount := get_service_account_name()
	serviceAccount != ""
	serviceAccount != "default"
	not service_account_exists(serviceAccount)
	msg := sprintf("serviceAccount %q not found", [serviceAccount])
}

get_service_account_name() := sa {
	sa := input.spec.template.spec.serviceAccountName
}

get_service_account_name() := sa {
	sa := input.spec.template.spec.deprecatedServiceAccount
}

get_service_account_name() := "" {
	not input.spec.template.spec.serviceAccountName
	not input.spec.template.spec.deprecatedServiceAccount
}

service_account_exists(serviceAccount) {
	some sa in data.objects
	sa.kind == "ServiceAccount"
	sa.metadata.namespace == input.metadata.namespace
	sa.metadata.name == serviceAccount
}