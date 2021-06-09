package replicas

import (
	"fmt"

	"golang.stackrox.io/kube-linter/internal/stringutils"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/replicas/internal/params"
)

const (
	templateKey = "minimum-replicas"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Minimum replicas",
		Key:         templateKey,
		Description: "Flag applications running fewer than the specified number of replicas",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				replicas, found := extract.Replicas(object.K8sObject)
				if !found {
					return nil
				}
				if int(replicas) >= p.MinReplicas {
					return nil
				}
				return []diagnostic.Diagnostic{
					{Message: fmt.Sprintf("object has %d %s but minimum required replicas is %d",
						replicas, stringutils.Ternary(replicas > 1, "replicas", "replica"),
						p.MinReplicas)},
				}
			}, nil
		}),
	})
}
