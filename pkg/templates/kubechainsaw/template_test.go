package kubechainsaw

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/templates/kubechainsaw/internal/params"
	v1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type fakeLintContext struct {
	objects []lintcontext.Object
}

func (f *fakeLintContext) Objects() []lintcontext.Object              { return f.objects }
func (f *fakeLintContext) InvalidObjects() []lintcontext.InvalidObject { return nil }

func TestAnalyzeDetectsWildcardVerbs(t *testing.T) {
	cr := &rbacV1.ClusterRole{
		ObjectMeta: metaV1.ObjectMeta{Name: "dangerous-role"},
		Rules: []rbacV1.PolicyRule{
			{APIGroups: []string{""}, Resources: []string{"secrets"}, Verbs: []string{"*"}},
		},
	}
	cr.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"})

	ctx := &fakeLintContext{objects: []lintcontext.Object{
		{Metadata: lintcontext.ObjectMetadata{FilePath: "test.yaml"}, K8sObject: cr},
	}}

	p := params.Params{}
	checkFn, err := analyze(p)
	require.NoError(t, err)

	diags := checkFn(ctx, ctx.Objects()[0])
	assert.NotEmpty(t, diags, "should detect wildcard verbs")

	found := false
	for _, d := range diags {
		if d.Metadata != nil && d.Metadata[diagnostic.MetaKeyRuleID] == "KC-002" {
			found = true
			assert.NotEmpty(t, d.Severity)
			assert.NotEmpty(t, d.Metadata[diagnostic.MetaKeyFingerprint])
			assert.NotEmpty(t, d.Metadata[diagnostic.MetaKeyRemediation])
		}
	}
	assert.True(t, found, "should have KC-002 finding")
}

func TestAnalyzeEarlyReturnForIrrelevantKind(t *testing.T) {
	cm := &v1.ConfigMap{
		ObjectMeta: metaV1.ObjectMeta{Name: "test-cm", Namespace: "default"},
	}
	cm.SetGroupVersionKind(schema.GroupVersionKind{Version: "v1", Kind: "ConfigMap"})

	ctx := &fakeLintContext{objects: []lintcontext.Object{
		{Metadata: lintcontext.ObjectMetadata{FilePath: "cm.yaml"}, K8sObject: cm},
	}}

	p := params.Params{}
	checkFn, err := analyze(p)
	require.NoError(t, err)

	diags := checkFn(ctx, ctx.Objects()[0])
	assert.Empty(t, diags, "ConfigMap should never have findings")
}

func TestAnalyzeRuleExclusion(t *testing.T) {
	cr := &rbacV1.ClusterRole{
		ObjectMeta: metaV1.ObjectMeta{Name: "dangerous-role"},
		Rules: []rbacV1.PolicyRule{
			{APIGroups: []string{""}, Resources: []string{"*"}, Verbs: []string{"*"}},
		},
	}
	cr.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"})

	ctx := &fakeLintContext{objects: []lintcontext.Object{
		{Metadata: lintcontext.ObjectMetadata{FilePath: "test.yaml"}, K8sObject: cr},
	}}

	p := params.Params{ExcludeRules: []string{"KC-001", "KC-002"}}
	checkFn, err := analyze(p)
	require.NoError(t, err)

	diags := checkFn(ctx, ctx.Objects()[0])
	for _, d := range diags {
		if d.Metadata != nil {
			assert.NotEqual(t, "KC-001", d.Metadata[diagnostic.MetaKeyRuleID])
			assert.NotEqual(t, "KC-002", d.Metadata[diagnostic.MetaKeyRuleID])
		}
	}
}
