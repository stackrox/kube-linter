package lint

import (
	"bytes"
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/run"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func junitTestPod(namespace, name string) *corev1.Pod {
	return &corev1.Pod{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
		ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: name},
	}
}

func TestFormatJUnit(t *testing.T) {
	result := run.Result{
		Checks: []config.Check{
			{Name: "latest-tag", Remediation: "Pin the image tag."},
			{Name: "run-as-non-root"},
		},
		Reports: []diagnostic.WithContext{
			{
				Diagnostic:  diagnostic.Diagnostic{Message: `container "app" uses image "nginx:latest" & <latest> tag`},
				Check:       "latest-tag",
				Remediation: "Pin the image tag.",
				Object: lintcontext.Object{
					Metadata:  lintcontext.ObjectMetadata{FilePath: "deploy/pod.yaml"},
					K8sObject: junitTestPod("prod", "web"),
				},
			},
			{
				Diagnostic: diagnostic.Diagnostic{Message: `container "app" uses image "redis:latest"`},
				Check:      "latest-tag",
				Object: lintcontext.Object{
					Metadata:  lintcontext.ObjectMetadata{FilePath: "deploy/cache.yaml"},
					K8sObject: junitTestPod("", "cache"),
				},
			},
			{
				Diagnostic:  diagnostic.Diagnostic{Message: "failed to decode"},
				Check:       "failed-to-load-object",
				Remediation: "Confirm that the file is accessible and is valid k8s yaml.",
				Object: lintcontext.Object{
					Metadata:  lintcontext.ObjectMetadata{FilePath: "deploy/broken.yaml"},
					K8sObject: &unstructured.Unstructured{},
				},
			},
		},
		Summary: run.Summary{
			ChecksStatus: run.ChecksFailed,
			CheckEndTime: time.Date(2026, 10, 26, 8, 30, 0, 0, time.UTC),
		},
	}

	var buf bytes.Buffer
	require.NoError(t, formatLintJUnit(&buf, result))
	assert.True(t, strings.HasPrefix(buf.String(), xml.Header), "output should start with the XML header")

	var got junitTestSuites
	require.NoError(t, xml.Unmarshal(buf.Bytes(), &got), "output must be valid XML")

	assert.Equal(t, "kube-linter", got.Name)
	assert.Equal(t, 4, got.Tests)
	assert.Equal(t, 3, got.Failures)
	require.Len(t, got.Suites, 3)

	latest := got.Suites[0]
	assert.Equal(t, "latest-tag", latest.Name)
	assert.Equal(t, 2, latest.Tests)
	assert.Equal(t, 2, latest.Failures)
	assert.Equal(t, "2026-10-26T08:30:00", latest.Timestamp)
	require.Len(t, latest.TestCases, 2)
	assert.Equal(t, "prod/web /v1, Kind=Pod", latest.TestCases[0].Name)
	assert.Equal(t, "latest-tag", latest.TestCases[0].ClassName)
	assert.Equal(t, "deploy/pod.yaml", latest.TestCases[0].File)
	require.NotNil(t, latest.TestCases[0].Failure)
	assert.Equal(t, `container "app" uses image "nginx:latest" & <latest> tag`, latest.TestCases[0].Failure.Message)
	assert.Equal(t, "latest-tag", latest.TestCases[0].Failure.Type)
	assert.Contains(t, latest.TestCases[0].Failure.Text, "remediation: Pin the image tag.")
	assert.Contains(t, latest.TestCases[0].Failure.Text, "file: deploy/pod.yaml")
	assert.Equal(t, "<no namespace>/cache /v1, Kind=Pod", latest.TestCases[1].Name)

	passing := got.Suites[1]
	assert.Equal(t, "run-as-non-root", passing.Name)
	assert.Equal(t, 1, passing.Tests)
	assert.Equal(t, 0, passing.Failures)
	require.Len(t, passing.TestCases, 1)
	assert.Equal(t, "run-as-non-root", passing.TestCases[0].Name)
	assert.Nil(t, passing.TestCases[0].Failure)

	loadFailure := got.Suites[2]
	assert.Equal(t, "failed-to-load-object", loadFailure.Name)
	assert.Equal(t, 1, loadFailure.Failures)
	require.Len(t, loadFailure.TestCases, 1)
	assert.Equal(t, "deploy/broken.yaml", loadFailure.TestCases[0].File)
	assert.Equal(t, "deploy/broken.yaml", loadFailure.TestCases[0].Name)
}

func TestFormatJUnit_UniqueCaseNames(t *testing.T) {
	pod := lintcontext.Object{
		Metadata:  lintcontext.ObjectMetadata{FilePath: "deploy/pod.yaml"},
		K8sObject: junitTestPod("prod", "web"),
	}
	loadFailure := func(path string) diagnostic.WithContext {
		return diagnostic.WithContext{
			Diagnostic: diagnostic.Diagnostic{Message: "failed to decode"},
			Check:      "failed-to-load-object",
			Object: lintcontext.Object{
				Metadata:  lintcontext.ObjectMetadata{FilePath: path},
				K8sObject: &unstructured.Unstructured{},
			},
		}
	}
	result := run.Result{
		Checks: []config.Check{{Name: "latest-tag"}},
		Reports: []diagnostic.WithContext{
			{Diagnostic: diagnostic.Diagnostic{Message: `container "app" uses image "nginx:latest"`}, Check: "latest-tag", Object: pod},
			{Diagnostic: diagnostic.Diagnostic{Message: `container "sidecar" uses image "envoy:latest"`}, Check: "latest-tag", Object: pod},
			loadFailure("deploy/a.yaml"),
			loadFailure("deploy/b.yaml #2"),
			loadFailure("deploy/b.yaml"),
			loadFailure("deploy/b.yaml"),
		},
	}

	var buf bytes.Buffer
	require.NoError(t, formatLintJUnit(&buf, result))
	var got junitTestSuites
	require.NoError(t, xml.Unmarshal(buf.Bytes(), &got))
	require.Len(t, got.Suites, 2)

	var names []string
	for _, tc := range got.Suites[0].TestCases {
		names = append(names, tc.Name)
	}
	assert.Equal(t, []string{"prod/web /v1, Kind=Pod", "prod/web /v1, Kind=Pod #2"}, names)

	names = nil
	for _, tc := range got.Suites[1].TestCases {
		names = append(names, tc.Name)
	}
	assert.Equal(t, []string{"deploy/a.yaml", "deploy/b.yaml #2", "deploy/b.yaml", "deploy/b.yaml #3"}, names)
}

func TestFormatJUnit_NoChecks(t *testing.T) {
	var buf bytes.Buffer
	require.NoError(t, formatLintJUnit(&buf, run.Result{}))

	var got junitTestSuites
	require.NoError(t, xml.Unmarshal(buf.Bytes(), &got))
	assert.Equal(t, 0, got.Tests)
	assert.Equal(t, 0, got.Failures)
	assert.Empty(t, got.Suites)
}

func TestFormatJUnit_RejectsOtherData(t *testing.T) {
	assert.Error(t, formatLintJUnit(&bytes.Buffer{}, "not a result"))
}

func TestCommand_JUnitOutput(t *testing.T) {
	out := filepath.Join(t.TempDir(), "kube-linter.xml")
	cmd := createLintCommand("--format", "junit", "--output", out, "./testdata/valid-pod.yaml")
	require.NoError(t, cmd.Execute())

	content, err := os.ReadFile(filepath.Clean(out))
	require.NoError(t, err)
	var got junitTestSuites
	require.NoError(t, xml.Unmarshal(content, &got))
	assert.Equal(t, "kube-linter", got.Name)
	assert.NotEmpty(t, got.Suites, "every enabled check should get a test suite")
	assert.Equal(t, 0, got.Failures)
}
