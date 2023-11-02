package extract

import (
	ocpSecV1 "github.com/openshift/api/security/v1"
	"golang.stackrox.io/kube-linter/pkg/k8sutil"
)

// SCCallowPrivilegedContainer extracts allowPrivilegedContainer from the given object, if available.
func SCCallowPrivilegedContainer(obj k8sutil.Object) (bool, bool) {
	if scc, ok := obj.(*ocpSecV1.SecurityContextConstraints); ok {
		return scc.AllowPrivilegedContainer, true
	}
	return false, false
}
