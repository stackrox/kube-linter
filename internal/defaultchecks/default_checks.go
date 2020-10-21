package defaultchecks

import (
	"golang.stackrox.io/kube-linter/internal/set"
)

var (
	// List is the list of built-in checks that are enabled by default.
	List = set.NewFrozenStringSet(
		"dangling-service",
		"deprecated-service-account-field",
		"env-var-secret",
		"no-extensions-v1beta",
		"no-read-only-root-fs",
		"privileged-container",
		"run-as-non-root",
		"unset-cpu-requirements",
		"unset-memory-requirements",
	)
)
