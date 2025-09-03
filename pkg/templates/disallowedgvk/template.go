package disallowedgvk

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/matcher"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/disallowedgvk/internal/params"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Disallowed API Objects",
		Key:         "disallowed-api-obj",
		Description: "Flag disallowed API object kinds",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.Any},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			groupMatcher, err := matcher.ForString(p.Group)
			if err != nil {
				return nil, fmt.Errorf("invalid group: %w", err)
			}
			versionMatcher, err := matcher.ForString(p.Version)
			if err != nil {
				return nil, fmt.Errorf("invalid version: %w", err)
			}
			kindMatcher, err := matcher.ForString(p.Kind)
			if err != nil {
				return nil, fmt.Errorf("invalid kind: %w", err)
			}
			return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				gvk := extract.GVK(object.K8sObject)
				if groupMatcher(gvk.Group) && versionMatcher(gvk.Version) && kindMatcher(gvk.Kind) {
					return []diagnostic.Diagnostic{{Message: fmt.Sprintf("disallowed API object found: %s", gvk)}}
				}
				return nil
			}, nil
		}),
	})
}
