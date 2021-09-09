package pdb

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/pdb/internal/params"
	pdbV1 "k8s.io/api/policy/v1"
	//intstrutil "k8s.io/apimachinery/pkg/util/intstr"
)

const (
	templateKey = "pod-disruption-budget"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Pod Disruption Budget not specified",
		Key:         templateKey,
		Description: "Flag DeploymentLike objects having no \"Pod Disruption Budget\" set in case if the replica count is more than 1",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{
				objectkinds.DeploymentLike,
				objectkinds.PodDisruptionBudget},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {

			criteria, ifPercent, err := getIntOrPercentValueSafelyFromString(p.MinimumMaxUnavailableCriteria)

			if err != nil {
				return nil, errors.Wrapf(err, "")
			}

			return func(lintCtx lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {

				var results []diagnostic.Diagnostic

				replicas, _ := extract.Replicas(object.K8sObject)
				if int(replicas) > 1 {
					namespaceForDeploymentLike := object.K8sObject.GetNamespace()
					PDBFound := false
					maxUnavailable := false
					labels := extract.Labels(object.K8sObject)

					pdbMaxUnavailable, pdbSelectors := gatherPDBs(lintCtx, namespaceForDeploymentLike, int(replicas))

					if len(pdbSelectors) > 0 {
						for pdbName, pdbSel := range pdbSelectors {

							if isMatch(labels, pdbSel) {
								PDBFound = true

								if ifPercent {
									criteria = int(math.Ceil(float64(criteria) * (float64(int(replicas))) / 100))
								}

								if (pdbMaxUnavailable != nil) && (pdbMaxUnavailable[pdbName] > criteria) {
									maxUnavailable = true
								}
							}
						}

						if !PDBFound {
							return []diagnostic.Diagnostic{{Message: fmt.Sprintf("Object has %d replicas, but there are no respective \"Pod Disruption Budget\" match found in the namespace \"%s\".\n", replicas, namespaceForDeploymentLike)}}
						} else if PDBFound && !maxUnavailable {
							return []diagnostic.Diagnostic{{Message: fmt.Sprintf("Object has a \"Pod Disruption Budget\" match found in the namespace \"%s\", but either it's \"maxUnavailable\" field is not set or it's less than the required criteria: %d.\n", namespaceForDeploymentLike, criteria)}}
						}
					}
				}
				return results
			}, nil
		}),
	})
}

func gatherPDBs(lintCtx lintcontext.LintContext, namespace string, replicas int) (map[string]int, map[string]map[string]string) {
	var listOfPDBSelector = make(map[string]map[string]string)
	var listOfPDBMaxUnavailable = map[string]int{}

	for _, object := range lintCtx.Objects() {
		pdb, ok := object.K8sObject.(*pdbV1.PodDisruptionBudget)
		if ok {
			if pdb.Namespace == namespace {
				nameOfPDB := pdb.Name
				PDBSelectors := pdb.Spec.Selector.MatchLabels
				listOfPDBSelector[nameOfPDB] = PDBSelectors

				unavailable, ifPercent, err := getIntOrPercentValueSafelyFromString((pdb.Spec.MaxUnavailable).String())

				if err == nil {
					if ifPercent {
						maxUnavailable := int(math.Ceil(float64(unavailable) * (float64(replicas)) / 100))
						listOfPDBMaxUnavailable[nameOfPDB] = maxUnavailable
					} else {
						maxUnavailable := unavailable
						listOfPDBMaxUnavailable[nameOfPDB] = maxUnavailable
					}
				}
			}
		} else {
			continue
		}
	}
	return listOfPDBMaxUnavailable, listOfPDBSelector
}

func isMatch(labels, pdbSelectors map[string]string) bool {
	if len(pdbSelectors) > len(labels) {
		return false
	}
	for key, value := range pdbSelectors {
		if res, found := labels[key]; !found || res != value {
			return false
		}
	}
	return true
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
