package lintcontext

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

func TestParseObjects(t *testing.T) {
	tests := []struct {
		name        string
		yamlData    string
		expectError bool
		expectCount int
		expectKind  string
		expectName  string
	}{
		{
			name: "valid Pod",
			yamlData: `apiVersion: v1
kind: Pod
metadata:
  name: test-pod
  namespace: default
spec:
  containers:
  - name: nginx
    image: nginx:latest
    ports:
    - containerPort: 80`,
			expectError: false,
			expectCount: 1,
			expectKind:  "Pod",
			expectName:  "test-pod",
		},
		{
			name: "valid Service",
			yamlData: `apiVersion: v1
kind: Service
metadata:
  name: test-service
  namespace: default
spec:
  selector:
    app: nginx
  ports:
  - port: 80
    targetPort: 80
  type: ClusterIP`,
			expectError: false,
			expectCount: 1,
			expectKind:  "Service",
			expectName:  "test-service",
		},
		{
			name: "Tekton Task CRD",
			yamlData: `apiVersion: tekton.dev/v1
kind: Task
metadata:
  name: hello-world-task
  namespace: default
spec:
  description: A simple hello world task
  steps:
  - name: hello
    image: alpine:latest
    command:
    - echo
    args:
    - "Hello World!"`,
			expectError: false,
			expectCount: 1,
			expectKind:  "Task",
			expectName:  "hello-world-task",
		},
		{
			name: "List with multiple objects",
			yamlData: `apiVersion: v1
kind: List
metadata: {}
items:
- apiVersion: v1
  kind: Pod
  metadata:
    name: pod1
  spec:
    containers:
    - name: nginx
      image: nginx:latest
- apiVersion: v1
  kind: Service
  metadata:
    name: service1
  spec:
    selector:
      app: nginx
    ports:
    - port: 80`,
			expectError: false,
			expectCount: 2,
			expectKind:  "Pod", // First object
			expectName:  "pod1",
		},
		{
			name: "invalid YAML",
			yamlData: `apiVersion: v1
kind: Pod
metadata:
  name: test-pod
spec:
  invalidField: this-should-not-be-here
  containers:
  - name: nginx
    image: nginx:latest
    invalidContainerField: also-invalid`,
			expectError: false, // parseObjects doesn't validate schema, only structure
			expectCount: 1,
			expectKind:  "Pod",
			expectName:  "test-pod",
		},
		{
			name: "malformed YAML",
			yamlData: `apiVersion: v1
kind: Pod
metadata:
  name: test-pod
spec:
  containers:
  - name: nginx
    image: nginx:latest
    ports:
    - containerPort: "invalid-port-type"`, // string instead of int
			expectError: true, // Should fail due to type mismatch
			expectCount: 0,
			expectKind:  "",
			expectName:  "",
		},
		{
			name: "unknown Kubernetes resource type",
			yamlData: `apiVersion: example.com/v1
kind: CustomResource
metadata:
  name: test-custom
  namespace: default
spec:
  customField: value`,
			expectError: false,
			expectCount: 1,
			expectKind:  "CustomResource",
			expectName:  "test-custom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objects, err := parseObjects([]byte(tt.yamlData), nil)

			if tt.expectError {
				assert.Error(t, err, "Expected parseObjects to return an error")
				assert.Len(t, objects, tt.expectCount)
			} else {
				assert.NoError(t, err, "Expected parseObjects to succeed")
				require.Len(t, objects, tt.expectCount, "Expected specific number of objects")

				if tt.expectCount > 0 {
					// Check first object
					firstObj := objects[0]
					assert.Equal(t, tt.expectKind, firstObj.GetObjectKind().GroupVersionKind().Kind)
					assert.Equal(t, tt.expectName, firstObj.GetName())

					// Additional validation for Pod objects
					if tt.expectKind == "Pod" {
						pod, ok := firstObj.(*corev1.Pod)
						require.True(t, ok, "Expected object to be a Pod")
						assert.Equal(t, "v1", pod.APIVersion)
						assert.Equal(t, "Pod", pod.Kind)
						assert.NotEmpty(t, pod.Spec.Containers, "Expected Pod to have containers")
					}
				}
			}
		})
	}
}

func TestParseObjectsWithCustomDecoder(t *testing.T) {
	// Test that parseObjects can handle CRDs by falling back to unstructured parsing
	tektonTaskYAML := `apiVersion: tekton.dev/v1
kind: Task
metadata:
  name: hello-world-task
spec:
  description: A simple hello world task
  steps:
  - name: hello
    image: alpine:latest
    command:
    - echo
    args:
    - "Hello World!"`

	// Test with default decoder (should succeed by falling back to unstructured)
	objects, err := parseObjects([]byte(tektonTaskYAML), nil)
	assert.NoError(t, err, "Expected Tekton Task to parse as unstructured with default decoder")
	assert.Len(t, objects, 1)
	assert.Equal(t, "Task", objects[0].GetObjectKind().GroupVersionKind().Kind)
	assert.Equal(t, "hello-world-task", objects[0].GetName())

	// Test with explicit decoder (should also succeed)
	objects, err = parseObjects([]byte(tektonTaskYAML), decoder)
	assert.NoError(t, err, "Expected Tekton Task to parse as unstructured with explicit decoder")
	assert.Len(t, objects, 1)
	assert.Equal(t, "Task", objects[0].GetObjectKind().GroupVersionKind().Kind)
	assert.Equal(t, "hello-world-task", objects[0].GetName())
}

func TestParseObjectsEmptyInput(t *testing.T) {
	// Test empty input
	objects, err := parseObjects([]byte(""), nil)
	assert.Error(t, err, "Expected empty input to return an error")
	assert.Empty(t, objects)

	// Test whitespace only
	objects, err = parseObjects([]byte("   \n  \t  \n"), nil)
	assert.Error(t, err, "Expected whitespace-only input to return an error")
	assert.Empty(t, objects)
}

func TestParseObjectsValidateObjectInterface(t *testing.T) {
	// Test that parsed objects implement the k8sutil.Object interface correctly
	podYAML := `apiVersion: v1
kind: Pod
metadata:
  name: test-pod
  namespace: test-namespace
  labels:
    app: test
  annotations:
    test: annotation
spec:
  containers:
  - name: nginx
    image: nginx:latest`

	objects, err := parseObjects([]byte(podYAML), nil)
	require.NoError(t, err)
	require.Len(t, objects, 1)

	pod := objects[0]

	// Test Object interface methods
	assert.Equal(t, "test-pod", pod.GetName())
	assert.Equal(t, "test-namespace", pod.GetNamespace())
	assert.Equal(t, map[string]string{"app": "test"}, pod.GetLabels())
	assert.Equal(t, map[string]string{"test": "annotation"}, pod.GetAnnotations())

	// Test GroupVersionKind
	gvk := pod.GetObjectKind().GroupVersionKind()
	assert.Empty(t, gvk.Group)
	assert.Equal(t, "v1", gvk.Version)
	assert.Equal(t, "Pod", gvk.Kind)
}
