package extract

import (
	ocpSecV1 "github.com/openshift/api/security/v1"
	"golang.stackrox.io/kube-linter/pkg/k8sutil"
)

// SCC extracts allowPrivilegedContainer from the given object, if available.
func SCCallowPrivilegedContainer(obj k8sutil.Object) (bool, bool) {
	switch scc := obj.(type) {
	case *ocpSecV1.SecurityContextConstraints:
		return scc.AllowPrivilegedContainer, true
	}
	return false, false
}
