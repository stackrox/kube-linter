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
	clusterRoleBindingGVK = schema.GroupVersionKind{
		Group:   rbacV1.GroupName,
		Version: "v1",
		Kind:    ClusterRoleBinding,
	}
)

func init() {
	registerObjectKind(ClusterRoleBinding, matcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == clusterRoleBindingGVK
	}))
}
