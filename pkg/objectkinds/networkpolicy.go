package objectkinds

import (
	v1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// NetworkPolicy represents Kubernetes NetworkPolicy objects.
	NetworkPolicy = "NetworkPolicy"
)

var (
	networkpolicyGVK = v1.SchemeGroupVersion.WithKind("NetworkPolicy")
)

func init() {
	RegisterObjectKind(NetworkPolicy, MatcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == networkpolicyGVK
	}))
}

// GetNetworkPolicyAPIVersion returns networkpolicy's apiversion
func GetNetworkPolicyAPIVersion() string {
	return networkpolicyGVK.GroupVersion().String()
}
