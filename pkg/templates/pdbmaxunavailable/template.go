package pdbmaxunavailable

import (
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/pdbmaxunavailable/internal/params"
	pdbV1 "k8s.io/api/policy/v1"
)

const (
	templateKey           = "pdb-max-unavailable"
	maxUnavailableZeroMsg = "MaxUnavailable is set to 0"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "No pod disruptions allowed - maxUnavailable",
		Key:         templateKey,
		Description: "Flag PodDisruptionBudgets whose maxUnavailable value will always prevent pod disruptions.",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{
				objectkinds.PodDisruptionBudget},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			return func(lintCtx lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {

				// Get the PDB provided
				pdb, ok := object.K8sObject.(*pdbV1.PodDisruptionBudget)
				if !ok {
					return nil
				}

				// If MaxUnavailable is set to 0 (both int and string)
				// return as no other checks will uniquely apply
				if pdb.Spec.MaxUnavailable != nil && pdb.Spec.MaxUnavailable.IntVal == 0 {
					return []diagnostic.Diagnostic{{
						Message: maxUnavailableZeroMsg,
					},
					}
				}

				return []diagnostic.Diagnostic{}

			}, nil
		}),
	})
}
