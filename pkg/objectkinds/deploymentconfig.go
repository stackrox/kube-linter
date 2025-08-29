package objectkinds

import (
	ocsAppsV1 "github.com/openshift/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// DeploymentConfig represents OpenShift DeploymentConfig objects.
	DeploymentConfig = "DeploymentConfig"
)

var (
	deploymentConfigGVK = ocsAppsV1.SchemeGroupVersion.WithKind("DeploymentConfig")
)

func init() {
	RegisterObjectKind(DeploymentConfig, MatcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == deploymentConfigGVK
	}))
}
