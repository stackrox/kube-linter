package nodeaffinity

import (
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/nodeaffinity/internal/params"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Node Affinity",
		Key:         "no-node-affinity",
		Description: "Flag objects that don't have node affinity rules set",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(_ params.Params) (check.Func, error) {
			return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				podTemplateSpec, found := extract.PodTemplateSpec(object.K8sObject)
				if !found {
					return nil
				}
				if podTemplateSpec.Spec.Affinity == nil {
					return []diagnostic.Diagnostic{{Message: "object does not define any node affinity rules."}}
				}
				if podTemplateSpec.Spec.Affinity.NodeAffinity == nil {
					return []diagnostic.Diagnostic{{Message: "object does not define any node affinity rules."}}
				}
				return nil
			}, nil
		}),
	})
}
