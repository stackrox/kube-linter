package pdbminavailable

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/k8sutil"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/pdbminavailable/internal/params"
	pdbV1 "k8s.io/api/policy/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	templateKey = "pdb-min-available"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Pod Disruption Budget Configuration Checks",
		Key:         templateKey,
		Description: "Flag non-optimal PDB configurations",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{
				objectkinds.PodDisruptionBudget},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			return func(lintCtx lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {

				var results []diagnostic.Diagnostic

				// Get the PDB provided
				pdb, ok := object.K8sObject.(*pdbV1.PodDisruptionBudget)
				if !ok {
					return nil
				}

				// If MinAvailable isn't set, then no need to check
				if pdb.Spec.MinAvailable == nil {
					return results
				}

				// Extract the MinAvailable value from the spec and check if it's a number or percentage
				value, isPercent, err := getIntOrPercentValueSafelyFromString(pdb.Spec.MinAvailable.String())
				if err != nil {
					return []diagnostic.Diagnostic{
						{
							Message: fmt.Sprintf("pdb has invalid MinAvailable value: %v", err),
						},
					}
				}

				// If the value is a percentage, handle the case where the MinValue is set to 100%
				// as DeploymentLike's replica counts don't need to be compared
				if isPercent {
					if value == 100 {
						return []diagnostic.Diagnostic{
							{
								Message: "PDB has minimum available replicas set to 100 percent of replicas",
							},
						}
					}
				}

				// Build the label selector for the PDB to use for comparison
				labelSelector, err := metaV1.LabelSelectorAsSelector(&metaV1.LabelSelector{
					MatchLabels:      pdb.Spec.Selector.MatchLabels,
					MatchExpressions: pdb.Spec.Selector.MatchExpressions},
				)

				if err != nil {
					return []diagnostic.Diagnostic{
						{
							Message: fmt.Sprintf("pdb has invalid label selector: %v", err),
						},
					}
				}

				// Evaluate Deploymet Likes in the lintContext to see if they have MinAvailable set too low
				deploymentLikes, err := getDeploymentLikeObjects(lintCtx, labelSelector, pdb.Namespace)
				if err != nil {
					return []diagnostic.Diagnostic{
						{
							Message: fmt.Sprintf("Failed to retrieve deployments matching the PDB's label selector within namespace %s: %v", pdb.Namespace, err.Error()),
						},
					}
				}

				for _, dl := range deploymentLikes {
					pdbMinAvailable := value
					replicas, _ := extract.Replicas(dl)
					if isPercent {
						// Calulate the actual value of the MinAvailable with respect to the Replica count if a percentage is set
						pdbMinAvailable = int(math.Floor(float64(replicas) * (float64(value) / float64(100))))
					}
					if replicas <= int32(pdbMinAvailable) {
						results = append(results, diagnostic.Diagnostic{
							Message: fmt.Sprintf("Deployment %s has replicas less than or equal to the minimum available replicas set by its PDB.", dl.GetName()),
						})
					}
				}
				return results

			}, nil
		}),
	})
}

func getDeploymentLikeObjects(lintCtx lintcontext.LintContext, labelSelector labels.Selector, namespace string) ([]k8sutil.Object, error) {

	objectList := make([]k8sutil.Object, 0, len(lintCtx.Objects()))

	for _, obj := range lintCtx.Objects() {
		// Ensure that only DeploymentLike objects are processed
		if !objectkinds.IsDeploymentLike(obj.GetK8sObjectName().GroupVersionKind) {
			continue
		}

		// Ensure that only DeploymentLikes are in the same namespaces as the PDB
		if obj.GetK8sObjectName().Namespace != namespace {
			continue
		}

		// Build Deployment labelSelector
		// If there are no selectors on the object, then the PDB won't match the same pods as the Deployment Like
		selectors, exists := extract.Selector(obj.K8sObject)
		if exists {
			// If there are no Replicas set on the Deployment Like, it's not possible to compare to a PDB
			if _, exists := extract.Replicas(obj.K8sObject); exists {
				objLabelSelector, err := metaV1.LabelSelectorAsSelector(selectors)
				if err != nil {
					return objectList, err
				}
				// Find any Deployment Likes with with the same selector as the PDB
				if labelSelector.String() == objLabelSelector.String() {
					objectList = append(objectList, obj.K8sObject)
				}
			}
		}
	}

	return objectList, nil
}

func getIntOrPercentValueSafelyFromString(intOrStr string) (int, bool, error) {

	var IntOrPercentValue = string("int")
	intValue, ifString := strconv.Atoi(intOrStr)

	if ifString != nil {
		IntOrPercentValue = "percent"
	}

	switch IntOrPercentValue {
	case "int":
		return intValue, false, nil
	case "percent":
		isPercent := false
		s := intOrStr
		if strings.HasSuffix(s, "%") {
			isPercent = true
			s = strings.TrimSuffix(s, "%")
		} else {
			return 0, false, fmt.Errorf("invalid type: string %q is not a percentage", intOrStr)
		}
		v, err := strconv.Atoi(s)
		if err != nil {
			return 0, false, fmt.Errorf("invalid value %q: %v", intOrStr, err)
		}
		return v, isPercent, nil
	}
	return 0, false, fmt.Errorf("invalid type: string %q neither int nor percentage", intOrStr)
}
