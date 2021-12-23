package forbiddenannotation

import (
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/forbiddenannotation/internal/params"
	"golang.stackrox.io/kube-linter/pkg/templates/util"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Forbidden Annotation",
		Key:         "forbidden-annotation",
		Description: "Flag objects carrying at least one annotation matching the provided patterns",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.Any},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			return util.ConstructForbiddenMapMatcher(p.Key, p.Value, "annotation")
		}),
	})
}
