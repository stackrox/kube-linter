package objectkinds

import (
	rbacV1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// Clusterrolebinding represents Kubernetes ClusterRoleBinding objects. Case sensitive.
	Clusterrolebinding = "ClusterRoleBinding"
)

var (
	clusterrolebindingGVK = schema.GroupVersionKind{
		Group:   rbacV1.GroupName,
		Version: "v1",
		Kind:    Clusterrolebinding,
	}
)

func init() {
	registerObjectKind(Clusterrolebinding, matcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == clusterrolebindingGVK
	}))
}
