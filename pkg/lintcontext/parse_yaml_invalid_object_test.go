package lintcontext

import (
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
