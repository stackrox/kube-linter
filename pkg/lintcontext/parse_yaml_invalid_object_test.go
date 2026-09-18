package lintcontext

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsV1 "k8s.io/api/apps/v1"
)

// Regression test for #591: a manifest of a *known* kind that fails to decode must be
// reported as an invalid object instead of being silently downgraded to Unstructured.
// The downgrade makes the workload invisible to every pod-spec based check, so the user
// gets an unrelated diagnostic (dangling-service) and never learns about the real problem.
func TestParseObjectsKnownKindWithInvalidValueIsReported(t *testing.T) {
	const brokenDeployment = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  selector:
    matchLabels:
      name: my-label-value
  template:
    metadata:
      labels:
        name: my-label-value
    spec:
      containers:
        - name: my-container
          image: nginx:latest
          resources:
            limits:
              cpu: 10wrong_unit
`

	objects, err := parseObjects([]byte(brokenDeployment), nil)
	assert.Error(t, err, "a Deployment with an invalid quantity must not parse silently")
	assert.Empty(t, objects)
	if err != nil {
		assert.Contains(t, err.Error(), "quantities must match", "the error must name the real cause")
	}
}

// The fallback to unstructured exists for custom resources whose Go type we do not know.
// That case must keep working: only *registered* kinds are held to the typed decoder.
func TestParseObjectsUnknownKindStillFallsBackToUnstructured(t *testing.T) {
	const tektonTask = `apiVersion: tekton.dev/v1beta1
kind: Task
metadata:
  name: my-task
spec:
  steps:
    - name: echo
      image: alpine
`

	objects, err := parseObjects([]byte(tektonTask), nil)
	require.NoError(t, err, "custom resources must still parse as unstructured")
	require.Len(t, objects, 1)
	assert.Equal(t, "Task", objects[0].GetObjectKind().GroupVersionKind().Kind)
}

// A valid manifest of a known kind must still decode into its typed Go representation,
// otherwise the checks that rely on extract.PodTemplateSpec silently stop seeing it.
func TestParseObjectsValidKnownKindStaysTyped(t *testing.T) {
	const okDeployment = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  selector:
    matchLabels:
      name: my-label-value
  template:
    metadata:
      labels:
        name: my-label-value
    spec:
      containers:
        - name: my-container
          image: nginx:latest
          resources:
            limits:
              cpu: 10m
`

	objects, err := parseObjects([]byte(okDeployment), nil)
	require.NoError(t, err)
	require.Len(t, objects, 1)
	assert.IsType(t, &appsV1.Deployment{}, objects[0])
}

// Reporting unreadable manifests must not turn section separators into findings.
// A document that holds nothing but comments carries no manifest at all, so it is
// neither a valid object nor an invalid one.
func TestLoadObjectsFromReaderSkipsCommentOnlyDocuments(t *testing.T) {
	const doc = `---
# Secret test cases
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  selector:
    matchLabels:
      name: my-label-value
  template:
    metadata:
      labels:
        name: my-label-value
    spec:
      containers:
        - name: my-container
          image: nginx:latest
---
# ConfigMap test cases
---
`

	ctx := newCtx(Options{})
	require.NoError(t, ctx.loadObjectsFromReader("test.yaml", strings.NewReader(doc)))
	assert.Empty(t, ctx.InvalidObjects(), "comment-only documents must not be reported as unreadable")
	require.Len(t, ctx.Objects(), 1)
	assert.Equal(t, "my-deployment", ctx.Objects()[0].K8sObject.GetName())
}

// A document marker may carry a comment on the same line. The start marker is consumed
// by the YAML reader, but the end marker is not: "... # end of section" reaches the
// parser as the whole body of a document and used to be reported as unreadable.
func TestLoadObjectsFromReaderSkipsDocumentMarkersWithComments(t *testing.T) {
	const doc = `--- # deployment section
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  selector:
    matchLabels:
      name: my-label-value
  template:
    metadata:
      labels:
        name: my-label-value
    spec:
      containers:
        - name: my-container
          image: nginx:latest
--- # nothing below this one
... # end of section
`

	ctx := newCtx(Options{})
	require.NoError(t, ctx.loadObjectsFromReader("test.yaml", strings.NewReader(doc)))
	assert.Empty(t, ctx.InvalidObjects(), "a document marker with a trailing comment carries no manifest")
	require.Len(t, ctx.Objects(), 1)
	assert.Equal(t, "my-deployment", ctx.Objects()[0].K8sObject.GetName())
}

// The marker must not swallow real content that happens to start with one: a tagged
// document ("--- !SomeTag") has something to decode and has to reach the parser.
func TestIsBlankDocumentKeepsContentAfterMarker(t *testing.T) {
	for _, tc := range []struct {
		name  string
		doc   string
		blank bool
	}{
		{name: "bare start marker", doc: "---", blank: true},
		{name: "bare end marker", doc: "...", blank: true},
		{name: "start marker with comment", doc: "--- # section two", blank: true},
		{name: "end marker with comment", doc: "... # end of section", blank: true},
		{name: "end marker then comment line", doc: "...\n# trailing note", blank: true},
		{name: "marker followed by a tag", doc: "--- !SomeTag", blank: false},
		{name: "marker glued to content", doc: "---foo: bar", blank: false},
		{name: "plain manifest", doc: "kind: Deployment", blank: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.blank, isBlankDocument([]byte(tc.doc)))
		})
	}
}
