package kubelinter.teamplate.latesttag

import data.kubelinter.objectkinds.is_deployment_like
import future.keywords.in

deny contains msg if {
	is_deployment_like
	violation_messages := [msg |
		m := matches[_]
		msg := sprintf(
			"The container %q is using an invalid container image, %q. Please use images that are not blocked by the `BlockList` criteria : %q",
			[m.image, m.image, m.pattern],
		)
	]
	msg := violation_messages
}

deny contains msg if {
	is_deployment_like
	violation_messages := [msg |
		m := not_matches[_]
		msg := sprintf(
			"The container %q is using an invalid container image, %q. Please use images that are allowed by the `AllowList` criteria : %q",
			[m.image, m.image, m.pattern],
		)
	]
	msg := violation_messages
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
