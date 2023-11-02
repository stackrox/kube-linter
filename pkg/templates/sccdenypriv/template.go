package sccdenypriv

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/sccdenypriv/internal/params"
)

const (
	templateKey = "scc-deny-privileged-container"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "SecurityContextConstraints allowPrivilegedContainer",
		Key:         templateKey,
		Description: "Flag SCC with allowPrivilegedContainer set to true",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.SecurityContextConstraints},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				state, found := extract.SCCallowPrivilegedContainer(object.K8sObject)
				if found && state == p.AllowPrivilegedContainer {
					return []diagnostic.Diagnostic{
						{Message: fmt.Sprintf("SCC has allowPrivilegedContainer set to %v", state)},
					}
				}
				return nil
			}, nil
		}),
	})
}
