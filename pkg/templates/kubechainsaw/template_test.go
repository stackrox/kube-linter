package kubechainsaw

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	kcModels "github.com/ugiordan/kube-chainsaw/pkg/models"
	"golang.stackrox.io/kube-linter/pkg/checkregistry"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/run"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/kubechainsaw/internal/convert"
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

func TestTemplateRegistersConverterKinds(t *testing.T) {
	// The engine only invokes the check for objects matching the template's
	// SupportedObjectKinds, so the list must cover every kind the converter
	// handles and must not fall back to Any.
	tmpl, found := templates.Get("kube-chainsaw")
	require.True(t, found)
	require.NotContains(t, tmpl.SupportedObjectKinds.ObjectKinds, objectkinds.Any)

	matcher, err := objectkinds.ConstructMatcher(tmpl.SupportedObjectKinds.ObjectKinds...)
	require.NoError(t, err)

	for _, tt := range []struct {
		gvk  schema.GroupVersionKind
		want bool
	}{
		{schema.GroupVersionKind{Version: "v1", Kind: "ServiceAccount"}, true},
		{schema.GroupVersionKind{Version: "v1", Kind: "Pod"}, true},
		{schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"}, true},
		{schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "Role"}, true},
		{schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRoleBinding"}, true},
		{schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "RoleBinding"}, true},
		{schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}, true},
		{schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "DaemonSet"}, true},
		{schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "StatefulSet"}, true},
		{schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "ReplicaSet"}, true},
		{schema.GroupVersionKind{Group: "batch", Version: "v1", Kind: "Job"}, true},
		{schema.GroupVersionKind{Group: "batch", Version: "v1", Kind: "CronJob"}, true},
		{schema.GroupVersionKind{Version: "v1", Kind: "ConfigMap"}, false},
	} {
		assert.Equal(t, tt.want, matcher.Matches(tt.gvk), "unexpected match result for %s", tt.gvk)
	}
}

func TestAnalyzeDoesNotReuseFindingsAcrossContexts(t *testing.T) {
	newContext := func(verbs []string) *fakeLintContext {
		role := &rbacV1.ClusterRole{
			ObjectMeta: metaV1.ObjectMeta{Name: "test-role"},
			Rules: []rbacV1.PolicyRule{{
				APIGroups: []string{""},
				Resources: []string{"configmaps"},
				Verbs:     verbs,
			}},
		}
		role.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   "rbac.authorization.k8s.io",
			Version: "v1",
			Kind:    "ClusterRole",
		})
		return &fakeLintContext{objects: []lintcontext.Object{
			{K8sObject: role},
		}}
	}

	tests := map[string]struct {
		firstDangerous bool
	}{
		"safe then dangerous": {firstDangerous: false},
		"dangerous then safe": {firstDangerous: true},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			checkFn, err := analyze(params.Params{
				Rules: []string{"KC-002"}, // Wildcard verbs.
			})
			require.NoError(t, err)

			for _, dangerous := range []bool{
				tt.firstDangerous,
				!tt.firstDangerous,
			} {
				verbs := []string{"get"}
				if dangerous {
					verbs = []string{"*"}
				}
				ctx := newContext(verbs)
				diags := checkFn(ctx, ctx.Objects()[0])

				if dangerous {
					assert.NotEmpty(t, diags,
						"dangerous context must have findings")
				} else {
					assert.Empty(t, diags,
						"safe context must have no findings")
				}
			}
		})
	}
}

func TestAnalyzeOmittedNamespaceDefaults(t *testing.T) {
	// A pod without an explicit namespace runs in "default", so it must be
	// linked to bindings whose subjects target default/<sa>.
	pod := &v1.Pod{
		ObjectMeta: metaV1.ObjectMeta{Name: "pod"},
		Spec:       v1.PodSpec{ServiceAccountName: "sa"},
	}
	pod.SetGroupVersionKind(schema.GroupVersionKind{Version: "v1", Kind: "Pod"})
	binding := &rbacV1.ClusterRoleBinding{
		ObjectMeta: metaV1.ObjectMeta{Name: "binding"},
		RoleRef:    rbacV1.RoleRef{Kind: "ClusterRole", Name: "cluster-admin"},
		Subjects:   []rbacV1.Subject{{Kind: "ServiceAccount", Name: "sa", Namespace: "default"}},
	}
	binding.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRoleBinding"})

	ctx := &fakeLintContext{objects: []lintcontext.Object{{K8sObject: pod}, {K8sObject: binding}}}
	fn, err := analyze(params.Params{Rules: []string{"KC-013"}})
	require.NoError(t, err)

	require.NotEmpty(t, fn(ctx, ctx.objects[0]), "pod with omitted namespace should match default/sa binding")
}

func TestAnalyzeLegacyServiceAccount(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "input.yaml")
	require.NoError(t, os.WriteFile(path, []byte(`apiVersion: v1
kind: Pod
metadata:
  name: pod
  namespace: default
spec:
  serviceAccount: sa
  containers:
  - name: app
    image: nginx
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: admin
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: sa
  namespace: default
`), 0o600))

	contexts, err := lintcontext.CreateContexts(nil, path)
	require.NoError(t, err)
	require.Len(t, contexts, 1)

	fn, err := analyze(params.Params{Rules: []string{"KC-013"}})
	require.NoError(t, err)

	found := false
	for _, obj := range contexts[0].Objects() {
		if len(fn(contexts[0], obj)) > 0 {
			found = true
		}
	}
	require.True(t, found, "legacy serviceAccount must retain the privilege chain")
}

func TestAnalyzeDuplicateIdentity(t *testing.T) {
	dir := t.TempDir()
	for name, verb := range map[string]string{"a-dangerous.yaml": "*", "b-safe.yaml": "get"} {
		content := "apiVersion: rbac.authorization.k8s.io/v1\nkind: ClusterRole\nmetadata:\n  name: shared\nrules:\n- apiGroups: ['']\n  resources: [configmaps]\n  verbs: ['" + verb + "']\n"
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600))
	}

	contexts, err := lintcontext.CreateContexts(nil, dir)
	require.NoError(t, err)
	require.Len(t, contexts, 1)

	fn, err := analyze(params.Params{Rules: []string{"KC-002"}})
	require.NoError(t, err)

	// The first call reports the ambiguity instead of silently keeping the
	// last definition, so the dangerous manifest cannot be overwritten.
	var diags []diagnostic.Diagnostic
	for _, obj := range contexts[0].Objects() {
		diags = append(diags, fn(contexts[0], obj)...)
	}
	require.NotEmpty(t, diags, "duplicate identity must be reported")
	assert.Contains(t, diags[0].Message, "duplicate ClusterRole")
}

func TestAnalyzeSCCResourceNames(t *testing.T) {
	role := &rbacV1.ClusterRole{
		ObjectMeta: metaV1.ObjectMeta{Name: "scc-user"},
		Rules: []rbacV1.PolicyRule{{
			APIGroups:     []string{"security.openshift.io"},
			Resources:     []string{"securitycontextconstraints"},
			Verbs:         []string{"use"},
			ResourceNames: []string{"restricted-v2"},
		}},
	}
	role.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"})

	ctx := &fakeLintContext{objects: []lintcontext.Object{{K8sObject: role}}}
	fn, err := analyze(params.Params{})
	require.NoError(t, err)

	found := false
	for _, d := range fn(ctx, ctx.objects[0]) {
		if strings.Contains(d.Message, "SecurityContextConstraints") {
			found = true
			assert.Contains(t, d.Message, "restricted-v2",
				"restricted resourceNames must be preserved in the finding")
		}
	}
	assert.True(t, found, "expected an SCC finding")
}

func TestAnalyzeGroupSubjectBinding(t *testing.T) {
	// A service account receives grants through its implicit identities, so a
	// binding to the system:serviceaccounts:<ns> group must link the pod to
	// cluster-admin the same way a direct ServiceAccount subject does.
	pod := &v1.Pod{
		ObjectMeta: metaV1.ObjectMeta{Name: "pod", Namespace: "default"},
		Spec:       v1.PodSpec{ServiceAccountName: "sa"},
	}
	pod.SetGroupVersionKind(schema.GroupVersionKind{Version: "v1", Kind: "Pod"})
	binding := &rbacV1.ClusterRoleBinding{
		ObjectMeta: metaV1.ObjectMeta{Name: "binding"},
		RoleRef:    rbacV1.RoleRef{Kind: "ClusterRole", Name: "cluster-admin"},
		Subjects:   []rbacV1.Subject{{Kind: "Group", Name: "system:serviceaccounts:default"}},
	}
	binding.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRoleBinding"})

	ctx := &fakeLintContext{objects: []lintcontext.Object{{K8sObject: pod}, {K8sObject: binding}}}
	fn, err := analyze(params.Params{Rules: []string{"KC-013"}})
	require.NoError(t, err)

	diags := fn(ctx, ctx.objects[0])
	require.NotEmpty(t, diags, "group grant must link the pod to cluster-admin")
	assert.Equal(t, diagnostic.SeverityCritical, diags[0].Severity)
}

func TestAnalyzeRejectsUnsupportedRules(t *testing.T) {
	// These rules analyze resource kinds the converter does not load, so
	// selecting them could never produce findings.
	for name, rule := range map[string]string{
		"NetworkPolicy ingress":        "KC-017",
		"NetworkPolicy egress":         "KC-018",
		"SecurityContextConstraints":   "KC-024",
		"external Service":             "KC-025",
		"unencrypted Route":            "KC-026",
		"broad external route":         "KC-027",
		"credential Secret":            "KC-028",
		"service account token Secret": "KC-029",
	} {
		t.Run(name, func(t *testing.T) {
			_, err := analyze(params.Params{Rules: []string{rule}})
			require.Error(t, err, "rule "+rule+" must be rejected while its resource kinds are unsupported")
			assert.Contains(t, err.Error(), "does not load")
		})
	}
}

func TestAnalyzeAcceptsSupportedRule(t *testing.T) {
	_, err := analyze(params.Params{Rules: []string{"KC-013"}})
	require.NoError(t, err)
}

func TestConvertRejectsDuplicateBindings(t *testing.T) {
	for name, objects := range map[string][]lintcontext.Object{
		"RoleBinding omitted and default namespaces": {
			{K8sObject: &rbacV1.RoleBinding{ObjectMeta: metaV1.ObjectMeta{Name: "shared"}}},
			{K8sObject: &rbacV1.RoleBinding{ObjectMeta: metaV1.ObjectMeta{Name: "shared", Namespace: "default"}}},
		},
		"ClusterRoleBinding": {
			{K8sObject: &rbacV1.ClusterRoleBinding{ObjectMeta: metaV1.ObjectMeta{Name: "shared"}}},
			{K8sObject: &rbacV1.ClusterRoleBinding{ObjectMeta: metaV1.ObjectMeta{Name: "shared"}}},
		},
	} {
		t.Run(name, func(t *testing.T) {
			objects[0].Metadata.FilePath = "first.yaml"
			objects[1].Metadata.FilePath = "second.yaml"
			_, err := convert.FromLintContext(&fakeLintContext{objects: objects})
			require.ErrorContains(t, err, "duplicate")
			require.ErrorContains(t, err, "first.yaml")
			require.ErrorContains(t, err, "second.yaml")
		})
	}
}

func TestConvertAllowsBindingsInDifferentNamespaces(t *testing.T) {
	ctx := &fakeLintContext{objects: []lintcontext.Object{
		{K8sObject: &rbacV1.RoleBinding{ObjectMeta: metaV1.ObjectMeta{Name: "shared", Namespace: "one"}}},
		{K8sObject: &rbacV1.RoleBinding{ObjectMeta: metaV1.ObjectMeta{Name: "shared", Namespace: "two"}}},
	}}
	resources, err := convert.FromLintContext(ctx)
	require.NoError(t, err)
	require.Len(t, resources.RoleBindings, 2)
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
		assert.NotEqual(t, diagnostic.SeverityInfo, severity, "info severity should be filtered out")
		assert.NotEqual(t, diagnostic.SeverityWarning, severity, "warning severity should be filtered out")
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
	// A missing suppressions file is a configuration error and must fail
	// instantiation instead of failing every lint run.
	_, err := analyze(params.Params{SuppressionsFile: "/nonexistent/file.yaml"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "/nonexistent/file.yaml")
}

func TestAnalyzeValidateCustomError(t *testing.T) {
	// Test the ValidateCustom error path
	p := params.Params{Rules: []string{"KC-INVALID"}}
	checkFn, err := analyze(p)
	require.Error(t, err)
	assert.Nil(t, checkFn)
	assert.Contains(t, err.Error(), "unknown rule ID")
}

func TestAnalyzeSuppressionSuccess(t *testing.T) {
	// Test the suppression success path
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

func TestFilterBySeverity(t *testing.T) {
	findings := []kcModels.Finding{
		{RuleID: "KC-001", Severity: kcModels.SeverityInfo},
		{RuleID: "KC-002", Severity: kcModels.SeverityCritical},
	}

	result := filterBySeverity(findings, "warning")
	require.Len(t, result, 1)
	assert.Equal(t, "KC-002", result[0].RuleID)

	result = filterBySeverity(findings, "")
	assert.Len(t, result, 2, "empty minSeverity should keep all findings")
}

func TestFilterBySeverityWithInvalidSeverity(t *testing.T) {
	// Test filterBySeverity error path
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
	// Test findingsForObject skipping suppressed findings
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
	cr := &rbacV1.ClusterRole{
		ObjectMeta: metaV1.ObjectMeta{Name: "test-role"},
		Rules:      []rbacV1.PolicyRule{},
	}
	cr.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"})

	ctx := &fakeLintContext{objects: []lintcontext.Object{
		{Metadata: lintcontext.ObjectMetadata{FilePath: "bad.yaml"}, K8sObject: nil},
		{Metadata: lintcontext.ObjectMetadata{FilePath: "good.yaml"}, K8sObject: cr},
	}}

	p := params.Params{}
	checkFn, err := analyze(p)
	require.NoError(t, err)

	// The first object reports the conversion failure.
	diags := checkFn(ctx, ctx.Objects()[0])
	require.Len(t, diags, 1)
	assert.Contains(t, diags[0].Message, "conversion error")
	assert.Equal(t, diagnostic.SeverityWarning, diags[0].Severity)

	// Other objects in the same run do not repeat it.
	assert.Empty(t, checkFn(ctx, ctx.Objects()[1]), "conversion error should only be reported once per run")

	// A new run over the same context reports it again.
	diags = checkFn(ctx, ctx.Objects()[0])
	require.Len(t, diags, 1)
	assert.Contains(t, diags[0].Message, "conversion error")
}

func TestAnalyzeInitializationFailurePersistsAcrossRuns(t *testing.T) {
	for name, duplicate := range map[string]bool{
		"missing suppression file":          false,
		"duplicate role conversion failure": true,
	} {
		t.Run(name, func(t *testing.T) {
			role := &rbacV1.ClusterRole{
				TypeMeta:   metaV1.TypeMeta{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "ClusterRole"},
				ObjectMeta: metaV1.ObjectMeta{Name: "shared"},
			}
			ctx := &fakeLintContext{objects: []lintcontext.Object{
				{Metadata: lintcontext.ObjectMetadata{FilePath: "first.yaml"}, K8sObject: role},
			}}
			check := &config.Check{Name: "review", Template: "kube-chainsaw"}
			if duplicate {
				ctx.objects = append(ctx.objects, lintcontext.Object{
					Metadata: lintcontext.ObjectMetadata{FilePath: "second.yaml"}, K8sObject: role.DeepCopy(),
				})
			} else {
				check.Params = map[string]interface{}{"suppressionsFile": filepath.Join(t.TempDir(), "missing.yaml")}
			}
			registry := checkregistry.New()
			err := registry.Register(check)
			if !duplicate && err != nil {
				// Rejecting invalid suppression configuration during instantiation is also correct.
				require.ErrorContains(t, err, "missing.yaml")
				return
			}
			require.NoError(t, err)
			for attempt := 1; attempt <= 2; attempt++ {
				result, err := run.Run([]lintcontext.LintContext{ctx}, registry, []string{"review"})
				require.NoError(t, err)
				assert.Equal(t, run.ChecksFailed, result.Summary.ChecksStatus, "run %d", attempt)
				assert.Len(t, result.Reports, 1, "report the failure once per run, including run %d", attempt)
			}
		})
	}
}
