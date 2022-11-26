package defaultchecks

import (
	"golang.stackrox.io/kube-linter/internal/set"
)

var (
	// List is the list of built-in checks that are enabled by default.
	List = set.NewFrozenStringSet(
		"dangling-service",
		"deprecated-service-account-field",
		"docker-sock",
		"drop-net-raw-capability",
		"duplicate-env-var",
		"env-var-secret",
		"host-ipc",
		"host-network",
		"host-pid",
		"invalid-target-ports",
		"latest-tag",
		"mismatching-selector",
		"no-anti-affinity",
		"no-extensions-v1beta",
		"no-read-only-root-fs",
		"non-existent-service-account",
		"privilege-escalation-container",
		"privileged-container",
		"probe-port",
		"run-as-non-root",
		"sensitive-host-mounts",
		"ssh-port",
		"unsafe-sysctls",
		"unset-cpu-requirements",
		"unset-memory-requirements",
	)
)
