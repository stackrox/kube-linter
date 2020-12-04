package all

import (
	// Import all check templates.
	_ "golang.stackrox.io/kube-linter/internal/templates/antiaffinity"
	_ "golang.stackrox.io/kube-linter/internal/templates/containercapabilities"
	_ "golang.stackrox.io/kube-linter/internal/templates/cpurequirements"
	_ "golang.stackrox.io/kube-linter/internal/templates/danglingservice"
	_ "golang.stackrox.io/kube-linter/internal/templates/deprecatedserviceaccount"
	_ "golang.stackrox.io/kube-linter/internal/templates/disallowedgvk"
	_ "golang.stackrox.io/kube-linter/internal/templates/envvar"
	_ "golang.stackrox.io/kube-linter/internal/templates/livenessprobe"
	_ "golang.stackrox.io/kube-linter/internal/templates/memoryrequirements"
	_ "golang.stackrox.io/kube-linter/internal/templates/mismatchingselector"
	_ "golang.stackrox.io/kube-linter/internal/templates/nonexistentserviceaccount"
	_ "golang.stackrox.io/kube-linter/internal/templates/ports"
	_ "golang.stackrox.io/kube-linter/internal/templates/privileged"
	_ "golang.stackrox.io/kube-linter/internal/templates/readinessprobe"
	_ "golang.stackrox.io/kube-linter/internal/templates/readonlyrootfs"
	_ "golang.stackrox.io/kube-linter/internal/templates/requiredannotation"
	_ "golang.stackrox.io/kube-linter/internal/templates/requiredlabel"
	_ "golang.stackrox.io/kube-linter/internal/templates/runasnonroot"
	_ "golang.stackrox.io/kube-linter/internal/templates/serviceaccount"
	_ "golang.stackrox.io/kube-linter/internal/templates/writablehostmount"
)
