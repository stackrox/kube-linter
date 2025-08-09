package kubelinter.template.serviceaccount

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny contains msg if {
	is_deployment_like
	serviceAccount := data.serviceaccount.serviceAccount
	podSA := get_service_account_name()
	regex.match(serviceAccount, podSA)
	msg := sprintf("found matching serviceAccount (%q)", [podSA])
}

get_service_account_name := sa if {
	# Check if automountServiceAccountToken is explicitly set to false
	not input.spec.template.spec.automountServiceAccountToken == false

	# Get service account name from pod spec
	sa := input.spec.template.spec.serviceAccountName
}

get_service_account_name := sa if {
	# Check if automountServiceAccountToken is explicitly set to false
	not input.spec.template.spec.automountServiceAccountToken == false

	# Fall back to deprecated service account field
	sa := input.spec.template.spec.deprecatedServiceAccount
}

get_service_account_name := "default" if {
	# Default to "default" if not specified
	not input.spec.template.spec.serviceAccountName
	not input.spec.template.spec.deprecatedServiceAccount
}
