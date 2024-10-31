package pdbunhealthypodevictionpolicy

import (
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/pdbunhealthypodevictionpolicy/internal/params"
	pdbV1 "k8s.io/api/policy/v1"
)

const (
	templateKey = "pdb-unhealthy-pod-eviction-policy"
)

func init() {
	templates.Register(check.Template{
		HumanName:   ".spec.unhealthyPodEvictionPolicy in PDB is set to default",
		Key:         templateKey,
		Description: "Flag PodDisruptionBudget objects that do not explicitly set unhealthyPodEvictionPolicy.",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{
				objectkinds.PodDisruptionBudget},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(_ params.Params) (check.Func, error) {
			return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				pdb, ok := object.K8sObject.(*pdbV1.PodDisruptionBudget)
				if !ok {
					return nil
				}
				if pdb.Spec.UnhealthyPodEvictionPolicy == nil {
					return []diagnostic.Diagnostic{{Message: "unhealthyPodEvictionPolicy is not explicitly set"}}
				}
				return nil
			}, nil
		}),
	})
}
