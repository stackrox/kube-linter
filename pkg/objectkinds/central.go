package objectkinds

import (
	"fmt"

	stackroxV1Alpha1 "golang.stackrox.io/kube-linter/pkg/crds/stackrox/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// Central represents StackRox Central objects. Case sensitive.
	Central = "Central"
)

var (
	CentralV1Alpha1 = stackroxV1Alpha1.GroupVersion.WithKind(Central)
)

func isCentral(gvk schema.GroupVersionKind) bool {
	return gvk == CentralV1Alpha1
}

func init() {
	RegisterObjectKind(Central, MatcherFunc(isCentral))
}

// GetCentralAPIVersion returns Central's APIVersion
func GetCentralAPIVersion(version string) string {
	return fmt.Sprintf("%s/%s", CentralV1Alpha1.Group, version)
}
