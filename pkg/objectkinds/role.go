package objectkinds

import (
	rbacV1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// Role represents Kubernetes Role objects. Case sensitive.
	Role = "Role"
)

var (
	// roleGVK represents Kubernetes Role objects. Case sensitive.
	roleGVK = rbacV1.SchemeGroupVersion.WithKind("Role")
)

func init() {
	RegisterObjectKind(Role, MatcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == roleGVK
	}))
}

// GetRoleAPIVersion returns Role's APIVersion
func GetRoleAPIVersion() string {
	return roleGVK.GroupVersion().String()
}
