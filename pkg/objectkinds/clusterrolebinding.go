package objectkinds

import (
	rbacV1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// ClusterRoleBinding represents Kubernetes ClusterRoleBinding objects. Case sensitive.
	ClusterRoleBinding = "ClusterRoleBinding"
)

var (
	clusterRoleBindingGVK = rbacV1.SchemeGroupVersion.WithKind("ClusterRoleBinding")
)

func init() {
	RegisterObjectKind(ClusterRoleBinding, MatcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == clusterRoleBindingGVK
	}))
}

// GetClusterRoleBindingAPIVersion returns ClusterRoleBinding's APIVersion
func GetClusterRoleBindingAPIVersion() string {
	return clusterRoleBindingGVK.GroupVersion().String()
}
