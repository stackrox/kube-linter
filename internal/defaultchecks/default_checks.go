package defaultchecks

import (
	"golang.stackrox.io/kube-linter/internal/set"
)

var (
	// List is the list of built-in checks that are enabled by default.
	List = set.NewFrozenStringSet(
		"privileged-container",
		"env-var-secret",
		"no-v1beta-deployments",
		"no-read-only-root-fs",
		"run-as-non-root",
	)
)
