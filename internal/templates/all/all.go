package all

import (
	// Import all check templates.
	_ "golang.stackrox.io/kube-linter/internal/templates/cpurequirements"
	_ "golang.stackrox.io/kube-linter/internal/templates/danglingservice"
	_ "golang.stackrox.io/kube-linter/internal/templates/disallowedgvk"
	_ "golang.stackrox.io/kube-linter/internal/templates/envvar"
	_ "golang.stackrox.io/kube-linter/internal/templates/memoryrequirements"
	_ "golang.stackrox.io/kube-linter/internal/templates/privileged"
	_ "golang.stackrox.io/kube-linter/internal/templates/readonlyrootfs"
	_ "golang.stackrox.io/kube-linter/internal/templates/requiredlabel"
	_ "golang.stackrox.io/kube-linter/internal/templates/runasnonroot"
)
