package kubelinter.template.latesttag

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny contains msg if {
	is_deployment_like
	some m in matches
	msg := sprintf(
		concat("", [
			"The container %q is using an invalid container image, %q. ",
			"Please use images that are not blocked by the `BlockList` criteria : %q",
		]),
		[m.image, m.image, m.pattern],
	)
}

deny contains msg if {
	is_deployment_like
	some m in not_matches
	msg := sprintf(
		concat("", [
			"The container %q is using an invalid container image, %q. ",
			"Please use images that are allowed by the `AllowList` criteria : %q",
		]),
		[m.image, m.image, m.pattern],
	)
}

not_matches := [{"image": c.image, "pattern": a} |
	some a in data.latesttag.allowList
	some c in input.spec.template.spec.containers
	not regex.match(a, c.image)
]

matches := [{"image": c.image, "pattern": a} |
	some c in input.spec.template.spec.containers
	some a in data.latesttag.blockList
	regex.match(a, c.image)
]
