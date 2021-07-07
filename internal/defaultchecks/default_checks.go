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
		"env-var-secret",
		"host-ipc",
		"sensitive-host-mounts",
		"host-network",
		"host-pid",
		"latest-tag",
		"mismatching-selector",
		"no-anti-affinity",
		"no-extensions-v1beta",
		"no-read-only-root-fs",
		"non-existent-service-account",
		"ssh-port",
		"privilege-escalation-container",
		"privileged-container",
		"run-as-non-root",
		"unsafe-sysctls",
		"unset-cpu-requirements",
		"unset-memory-requirements",
	)
)
