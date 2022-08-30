package objectkinds

import (
	rbacV1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// RoleBinding represents Kubernetes RoleBinding objects. Case sensitive.
	RoleBinding = "RoleBinding"
)

var (
	roleBindingGVK = rbacV1.SchemeGroupVersion.WithKind("RoleBinding")
)

func init() {
	RegisterObjectKind(RoleBinding, MatcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == roleBindingGVK
	}))
}

// GetRoleBindingAPIVersion returns RoleBinding's APIVersion
func GetRoleBindingAPIVersion() string {
	return roleBindingGVK.GroupVersion().String()
}
