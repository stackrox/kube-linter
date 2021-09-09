package objectkinds

import (
	"fmt"

	pdbV1 "k8s.io/api/policy/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// PodDisruptionBudget represents Kubernetes PodDisruptionBudget objects. Case sensitive.
	PodDisruptionBudget = "PodDisruptionBudget"
)

var (
	podDisruptionBudgetGVK = pdbV1.SchemeGroupVersion.WithKind("PodDisruptionBudget")
)

func init() {
	registerObjectKind(PodDisruptionBudget, matcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == podDisruptionBudgetGVK
	}))
}

// PodDisruptionBudgetAPIVersion returns PodDisruptionBudget's APIVersion
func PodDisruptionBudgetAPIVersion() string {
	return fmt.Sprintf("%s/%s", podDisruptionBudgetGVK.Group, podDisruptionBudgetGVK.Version)
}
