package sysctl

import (
	"fmt"
	"strings"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/sysctl/internal/params"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Unsafe Sysctls",
		Key:         "unsafe-sysctls",
		Description: "Flag unsafe sysctls",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				podSpec, found := extract.PodSpec(object.K8sObject)
				if !found {
					return nil
				}
				var results []diagnostic.Diagnostic
				if podSpec.SecurityContext != nil && podSpec.SecurityContext.Sysctls != nil {
					for _, ctl := range podSpec.SecurityContext.Sysctls {
						for _, unsafeCtl := range p.UnsafeSysCtls {
							if strings.HasPrefix(ctl.Name, unsafeCtl) {
								results = append(results, diagnostic.Diagnostic{Message: fmt.Sprintf("resource specifies unsafe sysctl %q.", ctl.Name)})
							}
						}
					}
				}
				return results
			}, nil
		}),
	})
}
