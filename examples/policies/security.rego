package kubernetes.admission

# Deny pods that run as root
deny[msg] {
    input.kind == "Pod"
    not input.spec.securityContext.runAsNonRoot

    msg := "Pods must not run as root"
}

# Deny pods without resource limits
deny[msg] {
    input.kind == "Pod"
    container := input.spec.containers[_]
    not container.resources.limits.memory

    msg := sprintf("Container %v must have memory limits", [container.name])
}

deny[msg] {
    input.kind == "Pod"
    container := input.spec.containers[_]
    not container.resources.limits.cpu

    msg := sprintf("Container %v must have CPU limits", [container.name])
}

# Deny privileged containers
deny[msg] {
    input.kind == "Pod"
    container := input.spec.containers[_]
    container.securityContext.privileged

    msg := sprintf("Container %v must not run in privileged mode", [container.name])
}

# Deny containers with host mounts
deny[msg] {
    input.kind == "Pod"
    volume := input.spec.volumes[_]
    volume.hostPath

    msg := sprintf("Volume %v must not use hostPath", [volume.name])
}

# Allow non-Pod resources
allow {
    input.kind != "Pod"
}

# Allow pods that meet all security requirements
allow {
    input.kind == "Pod"
    input.spec.securityContext.runAsNonRoot
    all_containers_have_limits
    no_privileged_containers
    no_host_mounts
}

# Helper functions
all_containers_have_limits {
    container := input.spec.containers[_]
    container.resources.limits.memory
    container.resources.limits.cpu
}

no_privileged_containers {
    container := input.spec.containers[_]
    not container.securityContext.privileged
}

no_host_mounts {
    volume := input.spec.volumes[_]
    not volume.hostPath
}