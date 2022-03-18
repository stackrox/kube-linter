package extract

import (
	"golang.stackrox.io/kube-linter/pkg/k8sutil"
	autoscalingV2Beta1 "k8s.io/api/autoscaling/v2beta1"
)

// HPAMinReplicas extracts minReplicas from the given object, if available.
func HPAMinReplicas(obj k8sutil.Object) (int32, bool) {
	if hpa, isHPA := obj.(*autoscalingV2Beta1.HorizontalPodAutoscaler); isHPA {
		if hpa.Spec.MinReplicas != nil {
			return *hpa.Spec.MinReplicas, true
		}
		// If numReplicas is a `nil` pointer, then it defaults to 1.
		return 1, true
	}
	return 0, false
}
