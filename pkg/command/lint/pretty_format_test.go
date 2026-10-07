package lint

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/run"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func prettyReport(file, name, check, msg string) diagnostic.WithContext {
	u := &unstructured.Unstructured{}
	u.SetAPIVersion("apps/v1")
	u.SetKind("Deployment")
	u.SetNamespace("default")
	u.SetName(name)
	return diagnostic.WithContext{
		Diagnostic:  diagnostic.Diagnostic{Message: msg},
		Check:       check,
		Remediation: "do the right thing",
		Object:      lintcontext.Object{Metadata: lintcontext.ObjectMetadata{FilePath: file}, K8sObject: u},
	}
}

func renderPretty(t *testing.T, res run.Result) string {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("NO_UNICODE", "1")
	t.Setenv("KLP_FIX", "")
	var buf bytes.Buffer
	require.NoError(t, formatLintPretty(&buf, res))
	return buf.String()
}

func TestPrettyNoFindings(t *testing.T) {
	out := renderPretty(t, run.Result{})
	assert.Contains(t, out, "no findings")
}

func TestPrettySingleObject(t *testing.T) {
	out := renderPretty(t, run.Result{Reports: []diagnostic.WithContext{
		prettyReport("insecure.yaml", "web", "drop-net-raw-capability", `container "a" does not drop NET_RAW capability`),
		prettyReport("insecure.yaml", "web", "drop-net-raw-capability", `container "b" has ADD capability: "NET_RAW"`),
		prettyReport("insecure.yaml", "web", "privileged-container", `container "a" is privileged`),
		prettyReport("insecure.yaml", "web", "latest-tag", `The container "a" is using an invalid container image, "nginx:latest".`),
	}})
	assert.Contains(t, out, "4 findings")
	assert.Contains(t, out, "NET_RAW capability configuration (2)")
	assert.Contains(t, out, "Privileged container")
	assert.Contains(t, out, `Container "b" adds the forbidden NET_RAW capability.`)
	assert.Contains(t, out, "Deployment apps/v1 in default, from insecure.yaml")
	assert.Contains(t, out, "exit code 1")
	assert.NotContains(t, out, "\x1b")
	assert.Less(t, strings.Index(out, "Security"), strings.Index(out, "Hygiene"))
}

func TestPrettyMultipleObjects(t *testing.T) {
	out := renderPretty(t, run.Result{Reports: []diagnostic.WithContext{
		prettyReport("a.yaml", "web", "host-network", "resource shares host's network namespace"),
		prettyReport("b.yaml", "legacy-redis", "no-liveness-probe", `container "redis" does not specify a liveness probe`),
	}})
	assert.Contains(t, out, "2 objects in 2 files")
	assert.Contains(t, out, "legacy-redis")
}

func TestPrettyFix(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("KLP_FIX", "1")
	var buf bytes.Buffer
	require.NoError(t, formatLintPretty(&buf, run.Result{Reports: []diagnostic.WithContext{
		prettyReport("a.yaml", "web", "host-pid", "object shares the host's process namespace"),
	}}))
	assert.Contains(t, buf.String(), "Fix: do the right thing")
}

func TestPrettyForceColor(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	t.Setenv("FORCE_COLOR", "1")
	var buf bytes.Buffer
	require.NoError(t, formatLintPretty(&buf, run.Result{Reports: []diagnostic.WithContext{
		prettyReport("a.yaml", "web", "host-pid", "x"),
	}}))
	assert.Contains(t, buf.String(), "\x1b[")
}

func TestPrettyWrapTruncates(t *testing.T) {
	s := prettyStyle{ell: "..."}
	lines := s.wrap(strings.Repeat("word ", 50), 20, 2)
	require.Len(t, lines, 2)
	assert.True(t, strings.HasSuffix(lines[1], "..."))
	assert.LessOrEqual(t, len(lines[1]), 20)
}
