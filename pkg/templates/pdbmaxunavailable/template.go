package pdbmaxunavailable

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/pdbmaxunavailable/internal/params"
	pdbV1 "k8s.io/api/policy/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	templateKey = "pdb-max-unavailable"
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

				if pdb.Spec.MaxUnavailable == nil {
					return []diagnostic.Diagnostic{}
				}

				maxUnavailable, err := intstr.GetScaledValueFromIntOrPercent(pdb.Spec.MaxUnavailable, 100, false)
				if err != nil {
					return []diagnostic.Diagnostic{{
						Message: fmt.Sprintf("maxUnavailable has invalid value [%s]", pdb.Spec.MaxUnavailable),
					}}
				}

				if maxUnavailable == 0 {
					return []diagnostic.Diagnostic{{
						Message: "MaxUnavailable is set to 0",
					}}
				}

				return []diagnostic.Diagnostic{}
			}, nil
		}),
	})
}
