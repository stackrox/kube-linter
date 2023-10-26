package objectkinds

import (
	ocpSecV1 "github.com/openshift/api/security/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// Service represents Kubernetes Service objects.
	SecurityContextConstraints = "SecurityContextConstraints"
)

var (
	sccGVK = ocpSecV1.SchemeGroupVersion.WithKind("SecurityContextConstraints")
)

func init() {
	RegisterObjectKind(SecurityContextConstraints, MatcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == sccGVK
	}))
}

// GetSCCAPIVersion returns SCC's apiversion
func GetSCCAPIVersion() string {
	return sccGVK.GroupVersion().String()
}
