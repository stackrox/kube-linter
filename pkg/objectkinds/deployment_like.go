package objectkinds

import (
	"fmt"

	appsV1 "k8s.io/api/apps/v1"
	batchV1 "k8s.io/api/batch/v1"
	batchV1Beta1 "k8s.io/api/batch/v1beta1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	deploymentLikeGroupKinds = func() map[schema.GroupKind]struct{} {
		m := make(map[schema.GroupKind]struct{})
		for _, gk := range []schema.GroupKind{
			{Group: appsV1.GroupName, Kind: "Deployment"},
			{Group: appsV1.GroupName, Kind: "DaemonSet"},
			{Group: appsV1.GroupName, Kind: "StatefulSet"},
			{Group: appsV1.GroupName, Kind: "ReplicaSet"},
			{Group: coreV1.GroupName, Kind: "Pod"},
			{Group: coreV1.GroupName, Kind: "ReplicationController"},
			{Group: batchV1.GroupName, Kind: "Job"},
			{Group: batchV1Beta1.GroupName, Kind: "CronJob"},
		} {
			if _, ok := m[gk]; ok {
				panic(fmt.Sprintf("group kind double-registered: %v", gk))
			}
			m[gk] = struct{}{}
		}
		return m
	}()
)

func isDeploymentLike(gvk schema.GroupVersionKind) bool {
	_, ok := deploymentLikeGroupKinds[gvk.GroupKind()]
	return ok
}

const (
	// DeploymentLike is the name of the DeploymentLike ObjectKind.
	DeploymentLike = "DeploymentLike"
)

func init() {
	registerObjectKind(DeploymentLike, matcherFunc(isDeploymentLike))
}
