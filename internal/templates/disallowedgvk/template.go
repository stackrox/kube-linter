package disallowedgvk

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/diagnostic"
	"golang.stackrox.io/kube-linter/internal/extract"
	"golang.stackrox.io/kube-linter/internal/lintcontext"
	"golang.stackrox.io/kube-linter/internal/matcher"
	"golang.stackrox.io/kube-linter/internal/objectkinds"
	"golang.stackrox.io/kube-linter/internal/templates"
	"golang.stackrox.io/kube-linter/internal/templates/disallowedgvk/internal/params"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Disallowed API Objects",
		Key:         "disallowed-api-obj",
		Description: "Flag disallowed API object kinds",
		SupportedObjectKinds: check.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.Any},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			groupMatcher, err := matcher.ForString(p.Group)
			if err != nil {
				return nil, errors.Wrap(err, "invalid group")
			}
			versionMatcher, err := matcher.ForString(p.Version)
			if err != nil {
				return nil, errors.Wrap(err, "invalid version")
			}
			kindMatcher, err := matcher.ForString(p.Kind)
			if err != nil {
				return nil, errors.Wrap(err, "invalid kind")
			}
			return func(_ *lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				gvk := extract.GVK(object.K8sObject)
				if groupMatcher(gvk.Group) && versionMatcher(gvk.Version) && kindMatcher(gvk.Kind) {
					return []diagnostic.Diagnostic{{Message: fmt.Sprintf("disallowed API object found: %s", gvk)}}
				}
				return nil
			}, nil
		}),
	})
}
