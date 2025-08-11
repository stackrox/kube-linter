package kubelinter.template.latesttag

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
	some a in input.params.latesttag.AllowList
	some c in input.object.spec.template.spec.containers
	not regex.match(a, c.image)
]

matches := [{"image": c.image, "pattern": a} |
	some a in input.params.latesttag.BlockList
	some c in input.object.spec.template.spec.containers
	regex.match(a, c.image)
]


# DeploymentLike matches various deployment-like objects
is_deployment_like if {
	input.object.kind == "Deployment"
	input.object.apiVersion == "apps/v1"
}

is_deployment_like if {
	input.object.kind == "DaemonSet"
	input.object.apiVersion == "apps/v1"
}

is_deployment_like if {
	input.object.kind == "StatefulSet"
	input.object.apiVersion == "apps/v1"
}

is_deployment_like if {
	input.object.kind == "ReplicaSet"
	input.object.apiVersion == "apps/v1"
}

is_deployment_like if {
	input.object.kind == "Pod"
	input.object.apiVersion == "v1"
}

is_deployment_like if {
	input.object.kind == "ReplicationController"
	input.object.apiVersion == "v1"
}

is_deployment_like if {
	input.object.kind == "Job"
	input.object.apiVersion == "batch/v1"
}

is_deployment_like if {
	input.object.kind == "CronJob"
	input.object.apiVersion == "batch/v1"
}

# OpenShift specific deployment-like objects
is_deployment_like if {
	input.object.kind == "DeploymentConfig"
	input.object.apiVersion == "apps.openshift.io/v1"
}
