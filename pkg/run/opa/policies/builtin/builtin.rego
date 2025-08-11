package pkg.run.opa.policies.builtin

is_ingress if {
	input.request.kind.kind == "Ingress"
	input.request.kind.group == "extensions"
	input.request.kind.version == "v1beta1"
}
