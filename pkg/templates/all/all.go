package all

import (
	// Import all check templates.
	_ "golang.stackrox.io/kube-linter/pkg/templates/antiaffinity"
	_ "golang.stackrox.io/kube-linter/pkg/templates/containercapabilities"
	_ "golang.stackrox.io/kube-linter/pkg/templates/cpurequirements"
	_ "golang.stackrox.io/kube-linter/pkg/templates/danglingservice"
	_ "golang.stackrox.io/kube-linter/pkg/templates/deprecatedserviceaccount"
	_ "golang.stackrox.io/kube-linter/pkg/templates/disallowedgvk"
	_ "golang.stackrox.io/kube-linter/pkg/templates/envvar"
	_ "golang.stackrox.io/kube-linter/pkg/templates/livenessprobe"
	_ "golang.stackrox.io/kube-linter/pkg/templates/memoryrequirements"
	_ "golang.stackrox.io/kube-linter/pkg/templates/mismatchingselector"
	_ "golang.stackrox.io/kube-linter/pkg/templates/nonexistentserviceaccount"
	_ "golang.stackrox.io/kube-linter/pkg/templates/ports"
	_ "golang.stackrox.io/kube-linter/pkg/templates/privileged"
	_ "golang.stackrox.io/kube-linter/pkg/templates/readinessprobe"
	_ "golang.stackrox.io/kube-linter/pkg/templates/readonlyrootfs"
	_ "golang.stackrox.io/kube-linter/pkg/templates/requiredannotation"
	_ "golang.stackrox.io/kube-linter/pkg/templates/requiredlabel"
	_ "golang.stackrox.io/kube-linter/pkg/templates/runasnonroot"
	_ "golang.stackrox.io/kube-linter/pkg/templates/serviceaccount"
	_ "golang.stackrox.io/kube-linter/pkg/templates/writablehostmount"
)
