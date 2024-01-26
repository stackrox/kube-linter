package objectkinds

import (
	"fmt"

	kedaV1Alpha1 "github.com/kedacore/keda/v2/apis/keda/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// ScaledObject represents Kubernetes ScaledObject objects. Case sensitive.
	ScaledObject = "ScaledObject"
)

var (
	ScaledObjectV1Alpha1 = kedaV1Alpha1.SchemeGroupVersion.WithKind(ScaledObject)
)

func isScaledObject(gvk schema.GroupVersionKind) bool {
	return gvk == ScaledObjectV1Alpha1
}

func init() {
	RegisterObjectKind(ScaledObject, MatcherFunc(isScaledObject))
}

// GetScaledObjectAPIVersion returns ScaledObject's APIVersion
func GetScaledObjectAPIVersion(version string) string {
	return fmt.Sprintf("%s/%s", ScaledObjectV1Alpha1.Group, version)
}
