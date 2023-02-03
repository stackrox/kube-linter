package objectkinds

import (
	policyv1 "k8s.io/api/policy/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// PodDisruptionBudget represents Kubernetes PodDisruptionBudget objects.
	PodDisruptionBudget = "PodDisruptionBudget"
)

var (
	pdbGVK = policyv1.SchemeGroupVersion.WithKind("PodDisruptionBudget")
)

func init() {
	RegisterObjectKind(PodDisruptionBudget, MatcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == pdbGVK
	}))
}

// GetPodDisruptionBudgetAPIVersion returns pdb's apiversion
func GetPodDisruptionBudgetAPIVersion() string {
	return pdbGVK.GroupVersion().String()
}
