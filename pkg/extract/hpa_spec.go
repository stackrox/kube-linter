package extract

import (
	kedaV1Alpha1 "github.com/kedacore/keda/v2/apis/keda/v1alpha1"
	"golang.stackrox.io/kube-linter/pkg/k8sutil"
	autoscalingV1 "k8s.io/api/autoscaling/v1"
	autoscalingV2 "k8s.io/api/autoscaling/v2"
	autoscalingV2Beta1 "k8s.io/api/autoscaling/v2beta1"
	autoscalingV2Beta2 "k8s.io/api/autoscaling/v2beta2"
)

// HPAMinReplicas extracts minReplicas from the given object, if available.
func HPAMinReplicas(obj k8sutil.Object) (int32, bool) {
	switch hpa := obj.(type) {
	case *autoscalingV2Beta1.HorizontalPodAutoscaler:
		return checkReplicas(hpa.Spec.MinReplicas)
	case *autoscalingV2Beta2.HorizontalPodAutoscaler:
		return checkReplicas(hpa.Spec.MinReplicas)
	case *autoscalingV2.HorizontalPodAutoscaler:
		return checkReplicas(hpa.Spec.MinReplicas)
	case *autoscalingV1.HorizontalPodAutoscaler:
		return checkReplicas(hpa.Spec.MinReplicas)
	case *kedaV1Alpha1.ScaledObject:
		return checkReplicas(hpa.Spec.MinReplicaCount)
	default:
		return 0, false
	}
}

func checkReplicas(minReplicas *int32) (int32, bool) {
	if minReplicas != nil {
		return *minReplicas, true
	}
	// If numReplicas is a `nil` pointer, then it defaults to 1.
	return 1, true
}

// HPAScaleTargetRefName extracts Spec.ScaleTargetRef.Name
func HPAScaleTargetRefName(obj k8sutil.Object) (string, bool) {
	switch hpa := obj.(type) {
	case *autoscalingV2Beta1.HorizontalPodAutoscaler:
		return hpa.Spec.ScaleTargetRef.Name, true
	case *autoscalingV2Beta2.HorizontalPodAutoscaler:
		return hpa.Spec.ScaleTargetRef.Name, true
	case *autoscalingV2.HorizontalPodAutoscaler:
		return hpa.Spec.ScaleTargetRef.Name, true
	case *autoscalingV1.HorizontalPodAutoscaler:
		return hpa.Spec.ScaleTargetRef.Name, true
	case *kedaV1Alpha1.ScaledObject:
		return hpa.Spec.ScaleTargetRef.Name, true
	default:
		return "", false
	}
}
