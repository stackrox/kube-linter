package objectkinds

import (
	"fmt"

	rbacV1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// ClusterRole represents Kubernetes ClusterRole objects. Case sensitive.
	ClusterRole = "ClusterRole"
)

var (
	clusterRoleGVK = rbacV1.SchemeGroupVersion.WithKind("ClusterRole")
)

func init() {
	RegisterObjectKind(ClusterRole, MatcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == clusterRoleGVK
	}))
}

// GetClusterRoleAPIVersion returns ClusterRole's APIVersion
func GetClusterRoleAPIVersion() string {
	return fmt.Sprintf("%s/%s", clusterRoleGVK.Group, clusterRoleGVK.Version)
}
