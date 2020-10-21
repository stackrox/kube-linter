package requiredannotation

import (
	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/objectkinds"
	"golang.stackrox.io/kube-linter/internal/templates"
	"golang.stackrox.io/kube-linter/internal/templates/requiredannotation/internal/params"
	"golang.stackrox.io/kube-linter/internal/templates/util"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Required Annotation",
		Key:         "required-annotation",
		Description: "Flag objects not carrying at least one annotation matching the provided patterns",
		SupportedObjectKinds: check.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.Any},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			return util.ConstructRequiredMapMatcher(p.Key, p.Value, "annotation")
		}),
	})
}
