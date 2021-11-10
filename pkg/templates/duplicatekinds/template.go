package duplicatekinds


import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/duplicatekinds/internal/params"
	v1 "k8s.io/api/core/v1"
)

type kindStruct struct{
	str string
    num int
}

/*{
    {"deployment", 0},
    {"daemonset", 0},
    {"statefulset", 0},
    {"service", 0},
}*/

func checkKindDuplicate(){

}

func init() {
	templates.Register(check.Template{
		HumanName:   "Duplicate Kind found",
		Key:         "duplicate-kinds",
		Description: "Flag when too many of a kind exist within a cluster",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.*},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(_ params.Params) (check.Func, error) {
			return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
                kind, found := object.K8sObject.(*v1.kind)
                if !found {
					return nil
				}
				var results []diagnostic.Diagnostic
				for _, duplicatekinds := range p.ForbiddenServiceTypes {
					if strings.EqualFold(string(service.Spec.Type), servicetype) {
						results = append(results, diagnostic.Diagnostic{Message: fmt.Sprintf("%q Duplicate Kind found.", duplicatekinds)})
					}
				}
				return results
			}, nil
		}),
	})
}
