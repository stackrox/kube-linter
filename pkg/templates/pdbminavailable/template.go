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
		HumanName:   "No pod disruptions allowed - minAvailable",
		Key:         templateKey,
		Description: "Flag PodDisruptionBudgets whose minAvailable value will always prevent pod disruptions.",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{
				objectkinds.PodDisruptionBudget},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			return minAvailableCheck, nil
		}),
	})
}

func minAvailableCheck(lintCtx lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {

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
				Message: fmt.Sprintf("PDB has invalid MinAvailable value: %v", err),
			},
		}
	}

	// If the value is a percentage, handle the case where the MinValue is set to 100%
	// as DeploymentLike's replica counts don't need to be compared
	if isPercent && value == 100 {
		return []diagnostic.Diagnostic{
			{
				Message: "PDB has minimum available replicas set to 100 percent of replicas",
			},
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
				Message: fmt.Sprintf("PDB has invalid label selector: %s", err),
			},
		}
	}

	// Evaluate Deploymet Likes in the lintContext to see if they have MinAvailable set too low
	deploymentLikes, err := getDeploymentLikeObjects(lintCtx, labelSelector, pdb.Namespace)
	if err != nil {
		return []diagnostic.Diagnostic{
			{
				Message: fmt.Sprintf("Failed to retrieve deployments matching the PDB's label selector within namespace %s: %s", pdb.Namespace, err),
			},
		}
	}

	for _, dl := range deploymentLikes {
		pdbMinAvailable := value
		replicas, _ := extract.Replicas(dl)
		if isPercent {
			// Calulate the actual value of the MinAvailable with respect to the Replica count if a percentage is set
			pdbMinAvailable = int(math.Ceil(float64(replicas) * (float64(value) / float64(100))))
		}
		if replicas <= int32(pdbMinAvailable) {
			results = append(results, diagnostic.Diagnostic{
				Message: fmt.Sprintf("The current number of replicas for deployment %s is equal to or lower than the minimum number of replicas specified by its PDB.", dl.GetName()),
			})
		}
	}
	return results
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
		if !exists {
			continue
		}
		// If there are no Replicas set on the Deployment Like, it's not possible to compare to a PDB
		if _, ok := extract.Replicas(obj.K8sObject); !ok {
			continue
		}

		objLabelSelector, err := metaV1.LabelSelectorAsSelector(selectors)
		if err != nil {
			return nil, err
		}
		// Find any Deployment Likes with the same selector as the PDB
		if labelSelector.String() == objLabelSelector.String() {
			objectList = append(objectList, obj.K8sObject)
		}
	}

	return objectList, nil
}

func getIntOrPercentValueSafelyFromString(intOrStr string) (int, bool, error) {
	intValue, ifString := strconv.Atoi(intOrStr)
	if ifString == nil {
		return intValue, false, nil
	}
	s := intOrStr
	if !strings.HasSuffix(s, "%") {
		return 0, false, fmt.Errorf("invalid type: string %q is not a percentage", intOrStr)
	}
	s = strings.TrimSuffix(s, "%")
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, false, fmt.Errorf("invalid value %q: %w", intOrStr, err)
	}
	return v, true, nil
}
