package kubelinter.template.disallowedgvk

import data.kubelinter.objectkinds.is_any

deny contains msg if {
	is_any
	regex.match(data.disallowedgvk.group, input.apiVersion)
	regex.match(data.disallowedgvk.version, input.apiVersion)
	regex.match(data.disallowedgvk.kind, input.kind)
	msg := sprintf(
		"disallowed API object found: %s/%s %s",
		[data.disallowedgvk.group, data.disallowedgvk.version, data.disallowedgvk.kind],
	)
}
