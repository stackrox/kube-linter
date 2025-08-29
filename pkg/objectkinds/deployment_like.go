package objectkinds

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	deploymentLikeGroupKinds = func() map[schema.GroupKind]struct{} {
		m := make(map[schema.GroupKind]struct{})
		for _, gk := range []schema.GroupKind{
			deploymentGVK.GroupKind(),
			daemonSetGVK.GroupKind(),
			deploymentConfigGVK.GroupKind(),
			statefulSetGVK.GroupKind(),
			replicaSetGVK.GroupKind(),
			podGVK.GroupKind(),
			replicationControllerGVK.GroupKind(),
			jobGVK.GroupKind(),
			cronJobGVK.GroupKind(),
		} {
			if _, ok := m[gk]; ok {
				panic(fmt.Sprintf("group kind double-registered: %v", gk))
			}
			m[gk] = struct{}{}
		}
		return m
	}()
)

func IsDeploymentLike(gvk schema.GroupVersionKind) bool {
	_, ok := deploymentLikeGroupKinds[gvk.GroupKind()]
	return ok
}

const (
	// DeploymentLike is the name of the DeploymentLike ObjectKind.
	DeploymentLike = "DeploymentLike"
)

func init() {
	RegisterObjectKind(DeploymentLike, MatcherFunc(IsDeploymentLike))
}
