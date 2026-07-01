package kubechainsaw

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	kcModels "github.com/ugiordan/kube-chainsaw/pkg/models"
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

func (f *fakeLintContext) Objects() []lintcontext.Object               { return f.objects }
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

func TestAnalyzeRuleInclusion(t *testing.T) {
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

	p := params.Params{Rules: []string{"KC-002"}}
	checkFn, err := analyze(p)
	require.NoError(t, err)

	diags := checkFn(ctx, ctx.Objects()[0])
	for _, d := range diags {
		if d.Metadata != nil {
			assert.Equal(t, "KC-002", d.Metadata[diagnostic.MetaKeyRuleID])
		}
	}
}

func TestAnalyzeSeverityFiltering(t *testing.T) {
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

	p := params.Params{MinSeverity: "high"}
	checkFn, err := analyze(p)
	require.NoError(t, err)

	diags := checkFn(ctx, ctx.Objects()[0])
	for _, d := range diags {
		severity := d.Severity
		assert.NotEqual(t, "info", severity, "info severity should be filtered out")
		assert.NotEqual(t, "warning", severity, "warning severity should be filtered out")
	}
}

func TestAnalyzeFindingsForObject_NonMatching(t *testing.T) {
	cr := &rbacV1.ClusterRole{
		ObjectMeta: metaV1.ObjectMeta{Name: "role-a"},
		Rules: []rbacV1.PolicyRule{
			{APIGroups: []string{""}, Resources: []string{"secrets"}, Verbs: []string{"*"}},
		},
	}
	cr.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"})

	cr2 := &rbacV1.ClusterRole{
		ObjectMeta: metaV1.ObjectMeta{Name: "role-b"},
		Rules:      []rbacV1.PolicyRule{},
	}
	cr2.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"})

	ctx := &fakeLintContext{objects: []lintcontext.Object{
		{Metadata: lintcontext.ObjectMetadata{FilePath: "test.yaml"}, K8sObject: cr},
		{Metadata: lintcontext.ObjectMetadata{FilePath: "test2.yaml"}, K8sObject: cr2},
	}}

	p := params.Params{}
	checkFn, err := analyze(p)
	require.NoError(t, err)

	checkFn(ctx, ctx.Objects()[0])

	diags := checkFn(ctx, ctx.Objects()[1])
	assert.Empty(t, diags, "role-b should have no findings")
}

func TestAnalyzeSuppressionFileError(t *testing.T) {
	cr := &rbacV1.ClusterRole{
		ObjectMeta: metaV1.ObjectMeta{Name: "test-role"},
		Rules:      []rbacV1.PolicyRule{},
	}
	cr.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"})

	ctx := &fakeLintContext{objects: []lintcontext.Object{
		{Metadata: lintcontext.ObjectMetadata{FilePath: "test.yaml"}, K8sObject: cr},
	}}

	p := params.Params{SuppressionsFile: "/nonexistent/file.yaml"}
	checkFn, err := analyze(p)
	require.NoError(t, err)

	diags1 := checkFn(ctx, ctx.Objects()[0])
	require.Len(t, diags1, 1)
	assert.Contains(t, diags1[0].Message, "failed to load suppressions")
	assert.Equal(t, "warning", diags1[0].Severity)

	diags2 := checkFn(ctx, ctx.Objects()[0])
	require.Len(t, diags2, 1)
	assert.Contains(t, diags2[0].Message, "failed to load suppressions")

	diags3 := checkFn(ctx, ctx.Objects()[0])
	assert.Empty(t, diags3, "suppression error should only be reported twice (once during analysis, once after)")
}

func TestAnalyzeValidateCustomError(t *testing.T) {
	// Test the ValidateCustom error path (template.go:35-38)
	p := params.Params{Rules: []string{"KC-INVALID"}}
	checkFn, err := analyze(p)
	require.Error(t, err)
	assert.Nil(t, checkFn)
	assert.Contains(t, err.Error(), "unknown rule ID")
}

func TestAnalyzeSuppressionSuccess(t *testing.T) {
	// Test the suppression success path (template.go:76)
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

	// Create a temp suppressions file
	tmpFile, err := os.CreateTemp("", "suppressions-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	suppressionYAML := `suppressions:
- rule_id: KC-002
  resource_name: dangerous-role
  reason: "accepted risk for test"
`
	_, err = tmpFile.WriteString(suppressionYAML)
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())

	p := params.Params{SuppressionsFile: tmpFile.Name()}
	checkFn, err := analyze(p)
	require.NoError(t, err)

	diags := checkFn(ctx, ctx.Objects()[0])
	// The finding should be suppressed
	for _, d := range diags {
		if d.Metadata != nil && d.Metadata[diagnostic.MetaKeyRuleID] == "KC-002" {
			t.Errorf("KC-002 finding should have been suppressed")
		}
	}
}

func TestFilterBySeverityWithInvalidSeverity(t *testing.T) {
	// Test filterBySeverity error path (template.go:175-183)
	// Create a fake finding
	findings := []kcModels.Finding{
		{RuleID: "KC-001", Severity: kcModels.SeverityHigh},
	}

	// Call filterBySeverity with an invalid severity string
	// This should return all findings unfiltered
	result := filterBySeverity(findings, "invalid-severity")
	assert.Equal(t, findings, result, "invalid severity should return all findings unfiltered")
}

func TestFindingsForObjectWithSuppressedFinding(t *testing.T) {
	// Test findingsForObject skipping suppressed findings (template.go:119-120)
	cr := &rbacV1.ClusterRole{
		ObjectMeta: metaV1.ObjectMeta{Name: "test-role"},
		Rules:      []rbacV1.PolicyRule{},
	}
	cr.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"})

	obj := lintcontext.Object{
		Metadata:  lintcontext.ObjectMetadata{FilePath: "test.yaml"},
		K8sObject: cr,
	}

	findings := []kcModels.Finding{
		{
			RuleID:            "KC-001",
			ResourceKind:      "ClusterRole",
			ResourceName:      "test-role",
			ResourceNamespace: "",
			Suppressed:        true,
		},
		{
			RuleID:            "KC-002",
			ResourceKind:      "ClusterRole",
			ResourceName:      "test-role",
			ResourceNamespace: "",
			Suppressed:        false,
		},
	}

	result := findingsForObject(findings, obj)
	require.Len(t, result, 1, "should only return non-suppressed finding")
	assert.Equal(t, "KC-002", result[0].RuleID)
}

func TestAnalyzeConversionError(t *testing.T) {
	ctx := &fakeLintContext{objects: []lintcontext.Object{
		{Metadata: lintcontext.ObjectMetadata{FilePath: "bad.yaml"}, K8sObject: nil},
	}}

	p := params.Params{}
	checkFn, err := analyze(p)
	require.NoError(t, err)

	// First call: error reported during analysis
	diags1 := checkFn(ctx, ctx.Objects()[0])
	require.Len(t, diags1, 1)
	assert.Contains(t, diags1[0].Message, "conversion error")
	assert.Equal(t, "warning", diags1[0].Severity)

	// Second call: error re-reported once via initErrReported flag
	diags2 := checkFn(ctx, ctx.Objects()[0])
	require.Len(t, diags2, 1)
	assert.Contains(t, diags2[0].Message, "conversion error")

	// Third call: no more reports
	diags3 := checkFn(ctx, ctx.Objects()[0])
	assert.Empty(t, diags3, "conversion error should stop after second report")
}
