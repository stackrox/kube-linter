package namespace

import (
	"strings"

	"golang.stackrox.io/kube-linter/internal/stringutils"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/namespace/internal/params"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Use Namespaces for Administrative Boundaries between Resources",
		Key:         "use-namespace",
		Description: "Flag resources with no namespace specified or using default namespace",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike, objectkinds.Service},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(_ params.Params) (check.Func, error) {
			return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				namespace := object.K8sObject.GetNamespace()
				ns := stringutils.OrDefault(namespace, "default")
				if strings.EqualFold(ns, "default") {
					return []diagnostic.Diagnostic{{Message: "object in default namespace"}}
				}
				return nil
			}, nil
		}),
	})
}
