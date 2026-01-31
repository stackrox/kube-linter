package objectkinds

import (
	"fmt"

	stackroxV1Alpha1 "golang.stackrox.io/kube-linter/pkg/crds/stackrox/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// SecuredCluster represents StackRox SecuredCluster objects. Case sensitive.
	SecuredCluster = "SecuredCluster"
)

var (
	SecuredClusterV1Alpha1 = stackroxV1Alpha1.GroupVersion.WithKind(SecuredCluster)
)

func isSecuredCluster(gvk schema.GroupVersionKind) bool {
	return gvk == SecuredClusterV1Alpha1
}

func init() {
	RegisterObjectKind(SecuredCluster, MatcherFunc(isSecuredCluster))
}

// GetSecuredClusterAPIVersion returns SecuredCluster's APIVersion
func GetSecuredClusterAPIVersion(version string) string {
	return fmt.Sprintf("%s/%s", SecuredClusterV1Alpha1.Group, version)
}
