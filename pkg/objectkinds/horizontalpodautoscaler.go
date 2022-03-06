package objectkinds

import (
	"fmt"

	autoscalingV2Beta1 "k8s.io/api/autoscaling/v2beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// HorizontalPodAutoscaler represents Kubernetes HorizontalPodAutoscaler objects. Case sensitive.
	HorizontalPodAutoscaler = "HorizontalPodAutoscaler"
)

var (
	horizontalPodAutoscalerGVK = autoscalingV2Beta1.SchemeGroupVersion.WithKind("HorizontalPodAutoscaler")
)

func init() {
	registerObjectKind(HorizontalPodAutoscaler, matcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == horizontalPodAutoscalerGVK
	}))
}

// GetHorizontalPodAutoscalerAPIVersion returns HorizontalPodAutoscaler's APIVersion
func GetHorizontalPodAutoscalerAPIVersion() string {
	return fmt.Sprintf("%s/%s", horizontalPodAutoscalerGVK.Group, horizontalPodAutoscalerGVK.Version)
}
