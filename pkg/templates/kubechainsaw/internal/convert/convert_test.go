package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
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
