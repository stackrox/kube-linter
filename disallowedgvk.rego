package kubelinter.template.disallowedgvk

import data.kubelinter.objectkinds.is_any

deny contains msg if {
	is_any
	group := data.disallowedgvk.group
	version := data.disallowedgvk.version
	kind := data.disallowedgvk.kind
	regex.match(group, input.apiVersion)
	regex.match(version, input.apiVersion)
	regex.match(kind, input.kind)
	msg := sprintf("disallowed API object found: %s/%s %s", [group, version, kind])
}
