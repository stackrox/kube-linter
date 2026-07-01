package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	appsV1 "k8s.io/api/apps/v1"
	batchV1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type fakeLintContext struct {
	objects []lintcontext.Object
}

func (f *fakeLintContext) Objects() []lintcontext.Object               { return f.objects }
func (f *fakeLintContext) InvalidObjects() []lintcontext.InvalidObject { return nil }

func TestConvertClusterRole(t *testing.T) {
	cr := &rbacV1.ClusterRole{
		ObjectMeta: metaV1.ObjectMeta{Name: "test-role"},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"secrets"},
				Verbs:     []string{"get", "list"},
			},
		},
	}
	cr.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"})

	ctx := &fakeLintContext{
		objects: []lintcontext.Object{
			{
				Metadata:  lintcontext.ObjectMetadata{FilePath: "rbac.yaml"},
				K8sObject: cr,
			},
		},
	}

	resources, err := FromLintContext(ctx)
	require.NoError(t, err)

	assert.Len(t, resources.ClusterRoles, 1)
	crData := resources.ClusterRoles["test-role"]
	require.NotNil(t, crData)
	assert.Equal(t, "rbac.yaml", crData.File)
	assert.Len(t, crData.Rules, 1)

	// Verify rules are map[string]interface{} format
	rule := crData.Rules[0]
	verbs, ok := rule["verbs"].([]interface{})
	require.True(t, ok, "verbs should be []interface{}")
	assert.Contains(t, verbs, "get")
	assert.Contains(t, verbs, "list")
}

func TestConvertRoleBinding(t *testing.T) {
	rb := &rbacV1.RoleBinding{
		ObjectMeta: metaV1.ObjectMeta{Name: "test-rb", Namespace: "default"},
		RoleRef:    rbacV1.RoleRef{Kind: "ClusterRole", Name: "admin"},
		Subjects: []rbacV1.Subject{
			{Kind: "ServiceAccount", Name: "my-sa", Namespace: "default"},
		},
	}
	rb.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "RoleBinding"})

	ctx := &fakeLintContext{objects: []lintcontext.Object{
		{Metadata: lintcontext.ObjectMetadata{FilePath: "rb.yaml"}, K8sObject: rb},
	}}

	resources, err := FromLintContext(ctx)
	require.NoError(t, err)
	require.Len(t, resources.RoleBindings, 1)

	bd := resources.RoleBindings[0]
	assert.Equal(t, "test-rb", bd.Name)
	assert.Equal(t, "default", bd.Namespace)
	assert.Equal(t, "ClusterRole", bd.RoleRef["kind"])
	assert.Equal(t, "admin", bd.RoleRef["name"])
	assert.Len(t, bd.Subjects, 1)
	assert.Equal(t, "ServiceAccount", bd.Subjects[0]["kind"])
}

func TestConvertPod(t *testing.T) {
	pod := &v1.Pod{
		ObjectMeta: metaV1.ObjectMeta{Name: "test-pod", Namespace: "default"},
		Spec:       v1.PodSpec{ServiceAccountName: "my-sa"},
	}
	pod.SetGroupVersionKind(schema.GroupVersionKind{Version: "v1", Kind: "Pod"})

	ctx := &fakeLintContext{objects: []lintcontext.Object{
		{Metadata: lintcontext.ObjectMetadata{FilePath: "pod.yaml"}, K8sObject: pod},
	}}

	resources, err := FromLintContext(ctx)
	require.NoError(t, err)
	require.Len(t, resources.Pods, 1)
	assert.Equal(t, "my-sa", resources.Pods["default/test-pod"].ServiceAccountName)
}

func TestConvertPodDefaultSA(t *testing.T) {
	pod := &v1.Pod{
		ObjectMeta: metaV1.ObjectMeta{Name: "test-pod", Namespace: "default"},
		Spec:       v1.PodSpec{},
	}
	pod.SetGroupVersionKind(schema.GroupVersionKind{Version: "v1", Kind: "Pod"})

	ctx := &fakeLintContext{objects: []lintcontext.Object{
		{Metadata: lintcontext.ObjectMetadata{FilePath: "pod.yaml"}, K8sObject: pod},
	}}

	resources, err := FromLintContext(ctx)
	require.NoError(t, err)
	assert.Equal(t, "default", resources.Pods["default/test-pod"].ServiceAccountName)
}

func TestConvertAggregatedClusterRole(t *testing.T) {
	cr := &rbacV1.ClusterRole{
		ObjectMeta: metaV1.ObjectMeta{Name: "agg-role"},
		AggregationRule: &rbacV1.AggregationRule{
			ClusterRoleSelectors: []metaV1.LabelSelector{
				{MatchLabels: map[string]string{"rbac": "true"}},
			},
		},
	}
	cr.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"})

	ctx := &fakeLintContext{objects: []lintcontext.Object{
		{Metadata: lintcontext.ObjectMetadata{FilePath: "agg.yaml"}, K8sObject: cr},
	}}

	resources, err := FromLintContext(ctx)
	require.NoError(t, err)

	crData := resources.ClusterRoles["agg-role"]
	require.NotNil(t, crData)
	_, hasAgg := crData.Doc["aggregationRule"]
	assert.True(t, hasAgg, "Doc should contain aggregationRule for KC-015 detection")
}

func TestConvertCronJob(t *testing.T) {
	cj := &batchV1.CronJob{
		ObjectMeta: metaV1.ObjectMeta{Name: "my-cj", Namespace: "default"},
		Spec: batchV1.CronJobSpec{
			JobTemplate: batchV1.JobTemplateSpec{
				Spec: batchV1.JobSpec{
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{ServiceAccountName: "cj-sa"},
					},
				},
			},
		},
	}
	cj.SetGroupVersionKind(schema.GroupVersionKind{Group: "batch", Version: "v1", Kind: "CronJob"})

	ctx := &fakeLintContext{objects: []lintcontext.Object{
		{Metadata: lintcontext.ObjectMetadata{FilePath: "cj.yaml"}, K8sObject: cj},
	}}

	resources, err := FromLintContext(ctx)
	require.NoError(t, err)
	require.Len(t, resources.Workloads, 1)
	assert.Equal(t, "cj-sa", resources.Workloads["CronJob/default/my-cj"].ServiceAccountName)
}

func TestConvertIgnoresUnknownTypes(t *testing.T) {
	cm := &v1.ConfigMap{
		ObjectMeta: metaV1.ObjectMeta{Name: "test-cm", Namespace: "default"},
	}
	cm.SetGroupVersionKind(schema.GroupVersionKind{Version: "v1", Kind: "ConfigMap"})

	ctx := &fakeLintContext{objects: []lintcontext.Object{
		{Metadata: lintcontext.ObjectMetadata{FilePath: "cm.yaml"}, K8sObject: cm},
	}}

	resources, err := FromLintContext(ctx)
	require.NoError(t, err)
	assert.True(t, resources.IsEmpty())
}

func TestConvertEmptyContext(t *testing.T) {
	ctx := &fakeLintContext{}
	resources, err := FromLintContext(ctx)
	require.NoError(t, err)
	assert.True(t, resources.IsEmpty())
}

func TestConvertDeployment(t *testing.T) {
	dep := &appsV1.Deployment{
		ObjectMeta: metaV1.ObjectMeta{Name: "test-dep", Namespace: "default"},
		Spec: appsV1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{ServiceAccountName: "dep-sa"},
			},
		},
	}
	dep.SetGroupVersionKind(schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"})

	ctx := &fakeLintContext{objects: []lintcontext.Object{
		{Metadata: lintcontext.ObjectMetadata{FilePath: "dep.yaml"}, K8sObject: dep},
	}}

	resources, err := FromLintContext(ctx)
	require.NoError(t, err)
	require.Len(t, resources.Workloads, 1)
	assert.Equal(t, "dep-sa", resources.Workloads["Deployment/default/test-dep"].ServiceAccountName)
	assert.Equal(t, "Deployment", resources.Workloads["Deployment/default/test-dep"].Kind)
}

func TestConvertDaemonSet(t *testing.T) {
	ds := &appsV1.DaemonSet{
		ObjectMeta: metaV1.ObjectMeta{Name: "test-ds", Namespace: "kube-system"},
		Spec: appsV1.DaemonSetSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{ServiceAccountName: "ds-sa"},
			},
		},
	}
	ds.SetGroupVersionKind(schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "DaemonSet"})

	ctx := &fakeLintContext{objects: []lintcontext.Object{
		{Metadata: lintcontext.ObjectMetadata{FilePath: "ds.yaml"}, K8sObject: ds},
	}}

	resources, err := FromLintContext(ctx)
	require.NoError(t, err)
	require.Len(t, resources.Workloads, 1)
	assert.Equal(t, "ds-sa", resources.Workloads["DaemonSet/kube-system/test-ds"].ServiceAccountName)
	assert.Equal(t, "DaemonSet", resources.Workloads["DaemonSet/kube-system/test-ds"].Kind)
}

func TestConvertStatefulSet(t *testing.T) {
	sts := &appsV1.StatefulSet{
		ObjectMeta: metaV1.ObjectMeta{Name: "test-sts", Namespace: "default"},
		Spec: appsV1.StatefulSetSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{ServiceAccountName: "sts-sa"},
			},
		},
	}
	sts.SetGroupVersionKind(schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "StatefulSet"})

	ctx := &fakeLintContext{objects: []lintcontext.Object{
		{Metadata: lintcontext.ObjectMetadata{FilePath: "sts.yaml"}, K8sObject: sts},
	}}

	resources, err := FromLintContext(ctx)
	require.NoError(t, err)
	require.Len(t, resources.Workloads, 1)
	assert.Equal(t, "sts-sa", resources.Workloads["StatefulSet/default/test-sts"].ServiceAccountName)
	assert.Equal(t, "StatefulSet", resources.Workloads["StatefulSet/default/test-sts"].Kind)
}

func TestConvertReplicaSet(t *testing.T) {
	rs := &appsV1.ReplicaSet{
		ObjectMeta: metaV1.ObjectMeta{Name: "test-rs", Namespace: "default"},
		Spec: appsV1.ReplicaSetSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{ServiceAccountName: "rs-sa"},
			},
		},
	}
	rs.SetGroupVersionKind(schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "ReplicaSet"})

	ctx := &fakeLintContext{objects: []lintcontext.Object{
		{Metadata: lintcontext.ObjectMetadata{FilePath: "rs.yaml"}, K8sObject: rs},
	}}

	resources, err := FromLintContext(ctx)
	require.NoError(t, err)
	require.Len(t, resources.Workloads, 1)
	assert.Equal(t, "rs-sa", resources.Workloads["ReplicaSet/default/test-rs"].ServiceAccountName)
	assert.Equal(t, "ReplicaSet", resources.Workloads["ReplicaSet/default/test-rs"].Kind)
}

func TestConvertJob(t *testing.T) {
	job := &batchV1.Job{
		ObjectMeta: metaV1.ObjectMeta{Name: "test-job", Namespace: "default"},
		Spec: batchV1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{ServiceAccountName: "job-sa"},
			},
		},
	}
	job.SetGroupVersionKind(schema.GroupVersionKind{Group: "batch", Version: "v1", Kind: "Job"})

	ctx := &fakeLintContext{objects: []lintcontext.Object{
		{Metadata: lintcontext.ObjectMetadata{FilePath: "job.yaml"}, K8sObject: job},
	}}

	resources, err := FromLintContext(ctx)
	require.NoError(t, err)
	require.Len(t, resources.Workloads, 1)
	assert.Equal(t, "job-sa", resources.Workloads["Job/default/test-job"].ServiceAccountName)
	assert.Equal(t, "Job", resources.Workloads["Job/default/test-job"].Kind)
}

func TestConvertServiceAccount(t *testing.T) {
	sa := &v1.ServiceAccount{
		ObjectMeta: metaV1.ObjectMeta{Name: "test-sa", Namespace: "default"},
	}
	sa.SetGroupVersionKind(schema.GroupVersionKind{Version: "v1", Kind: "ServiceAccount"})

	ctx := &fakeLintContext{objects: []lintcontext.Object{
		{Metadata: lintcontext.ObjectMetadata{FilePath: "sa.yaml"}, K8sObject: sa},
	}}

	resources, err := FromLintContext(ctx)
	require.NoError(t, err)
	require.Len(t, resources.ServiceAccounts, 1)
	assert.Equal(t, "test-sa", resources.ServiceAccounts["default/test-sa"].Name)
	assert.Equal(t, "default", resources.ServiceAccounts["default/test-sa"].Namespace)
}

func TestConvertClusterRoleBinding(t *testing.T) {
	crb := &rbacV1.ClusterRoleBinding{
		ObjectMeta: metaV1.ObjectMeta{Name: "test-crb"},
		RoleRef:    rbacV1.RoleRef{Kind: "ClusterRole", Name: "cluster-admin"},
		Subjects: []rbacV1.Subject{
			{Kind: "ServiceAccount", Name: "admin-sa", Namespace: "kube-system"},
		},
	}
	crb.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRoleBinding"})

	ctx := &fakeLintContext{objects: []lintcontext.Object{
		{Metadata: lintcontext.ObjectMetadata{FilePath: "crb.yaml"}, K8sObject: crb},
	}}

	resources, err := FromLintContext(ctx)
	require.NoError(t, err)
	require.Len(t, resources.ClusterRoleBindings, 1)

	bd := resources.ClusterRoleBindings[0]
	assert.Equal(t, "test-crb", bd.Name)
	assert.Empty(t, bd.Namespace)
	assert.Equal(t, "ClusterRole", bd.RoleRef["kind"])
	assert.Equal(t, "cluster-admin", bd.RoleRef["name"])
	assert.Len(t, bd.Subjects, 1)
}

func TestConvertRole(t *testing.T) {
	role := &rbacV1.Role{
		ObjectMeta: metaV1.ObjectMeta{Name: "test-role", Namespace: "default"},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"configmaps"},
				Verbs:     []string{"get", "list"},
			},
		},
	}
	role.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "Role"})

	ctx := &fakeLintContext{objects: []lintcontext.Object{
		{Metadata: lintcontext.ObjectMetadata{FilePath: "role.yaml"}, K8sObject: role},
	}}

	resources, err := FromLintContext(ctx)
	require.NoError(t, err)
	require.Len(t, resources.Roles, 1)

	roleData := resources.Roles["default/test-role"]
	require.NotNil(t, roleData)
	assert.Equal(t, "default", roleData.Namespace)
	assert.Len(t, roleData.Rules, 1)
}

func TestConvertMultipleObjects(t *testing.T) {
	cr := &rbacV1.ClusterRole{
		ObjectMeta: metaV1.ObjectMeta{Name: "test-cr"},
		Rules:      []rbacV1.PolicyRule{},
	}
	cr.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"})

	sa := &v1.ServiceAccount{
		ObjectMeta: metaV1.ObjectMeta{Name: "test-sa", Namespace: "default"},
	}
	sa.SetGroupVersionKind(schema.GroupVersionKind{Version: "v1", Kind: "ServiceAccount"})

	dep := &appsV1.Deployment{
		ObjectMeta: metaV1.ObjectMeta{Name: "test-dep", Namespace: "default"},
		Spec: appsV1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{ServiceAccountName: "test-sa"},
			},
		},
	}
	dep.SetGroupVersionKind(schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"})

	ctx := &fakeLintContext{objects: []lintcontext.Object{
		{Metadata: lintcontext.ObjectMetadata{FilePath: "manifest.yaml"}, K8sObject: cr},
		{Metadata: lintcontext.ObjectMetadata{FilePath: "manifest.yaml"}, K8sObject: sa},
		{Metadata: lintcontext.ObjectMetadata{FilePath: "manifest.yaml"}, K8sObject: dep},
	}}

	resources, err := FromLintContext(ctx)
	require.NoError(t, err)
	assert.Len(t, resources.ClusterRoles, 1)
	assert.Len(t, resources.ServiceAccounts, 1)
	assert.Len(t, resources.Workloads, 1)
}

func TestConvertWorkloadDefaultSA(t *testing.T) {
	// Test convertWorkload with empty ServiceAccountName (convert.go:134-135)
	dep := &appsV1.Deployment{
		ObjectMeta: metaV1.ObjectMeta{Name: "test-dep", Namespace: "default"},
		Spec: appsV1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{}, // Empty ServiceAccountName
			},
		},
	}
	dep.SetGroupVersionKind(schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"})

	ctx := &fakeLintContext{objects: []lintcontext.Object{
		{Metadata: lintcontext.ObjectMetadata{FilePath: "dep.yaml"}, K8sObject: dep},
	}}

	resources, err := FromLintContext(ctx)
	require.NoError(t, err)
	require.Len(t, resources.Workloads, 1)
	assert.Equal(t, "default", resources.Workloads["Deployment/default/test-dep"].ServiceAccountName)
}
