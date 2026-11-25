package extract

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestExtractPodTemplateSpec_WhenUnstructuredObject(t *testing.T) {
	tCases := []struct {
		name     string
		obj      *unstructured.Unstructured
		expected coreV1.PodTemplateSpec
		found    bool
	}{
		{
			name: "extract pod template spec from deployment like object",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "argoproj.io/v1alpha1",
					"kind":       "Rollout",
					"spec": map[string]interface{}{
						"template": map[string]interface{}{
							"metadata": map[string]interface{}{
								"name": "test",
							},
							"spec": map[string]interface{}{
								"containers": []interface{}{
									map[string]interface{}{
										"name":  "test-container",
										"image": "test-image",
									},
								},
							},
						},
					},
				},
			},
			expected: coreV1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: coreV1.PodSpec{
					Containers: []coreV1.Container{
						{
							Name:  "test-container",
							Image: "test-image",
						},
					},
				},
			},
			found: true,
		},
		{
			name: "extract pod template spec from not deployment like object",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "SomeTemplatingAPI",
					"spec": map[string]interface{}{
						"template": map[string]interface{}{
							"metadata": map[string]interface{}{
								"name": "test",
							},
							"spec": map[string]interface{}{
								"containers": []interface{}{
									map[string]interface{}{
										"name":  "test-container",
										"image": "test-image",
									},
								},
							},
						},
					},
				},
			},
			expected: coreV1.PodTemplateSpec{},
			found:    false,
		},
	}

	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			podTemplateSpec, found := PodTemplateSpec(tc.obj)
			assert.Equal(t, tc.found, found)
			assert.Equal(t, tc.expected, podTemplateSpec)
		})
	}
}

func TestExtractSelector_WhenUnstructuredObject(t *testing.T) {
	tCases := []struct {
		name     string
		obj      *unstructured.Unstructured
		expected *metav1.LabelSelector
		found    bool
	}{
		{
			name: "extract selector from deployment like object",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "argoproj.io/v1alpha1",
					"kind":       "Rollout",
					"spec": map[string]interface{}{
						"selector": map[string]interface{}{
							"matchLabels": map[string]interface{}{
								"app": "test",
							},
						},
						"template": map[string]interface{}{
							"metadata": map[string]interface{}{
								"name": "test",
							},
							"spec": map[string]interface{}{
								"containers": []interface{}{
									map[string]interface{}{
										"name":  "test-container",
										"image": "test-image",
									},
								},
							},
						},
					},
				},
			},
			expected: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "test",
				},
			},
			found: true,
		},
		{
			name: "extract selector from not deployment like object",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "SomeTemplatingAPI",
					"spec": map[string]interface{}{
						"selector": map[string]interface{}{
							"matchLabels": map[string]interface{}{
								"app": "test",
							},
						},
						"template": map[string]interface{}{
							"metadata": map[string]interface{}{
								"name": "test",
							},
							"spec": map[string]interface{}{
								"containers": []interface{}{
									map[string]interface{}{
										"name":  "test-container",
										"image": "test-image",
									},
								},
							},
						},
					},
				},
			},
			expected: nil,
			found:    false,
		},
	}

	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			selector, found := Selector(tc.obj)
			assert.Equal(t, tc.found, found)
			assert.Equal(t, tc.expected, selector)
		})
	}
}

func TestExtractReplicas_WhenUnstructuredObject(t *testing.T) {
	tCases := []struct {
		name     string
		obj      *unstructured.Unstructured
		expected int32
		found    bool
	}{
		{
			name: "extract replicas from deployment like object",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "argoproj.io/v1alpha1",
					"kind":       "Rollout",
					"spec": map[string]interface{}{
						"replicas": int64(3),
						"selector": map[string]interface{}{
							"matchLabels": map[string]interface{}{
								"app": "test",
							},
						},
						"template": map[string]interface{}{
							"metadata": map[string]interface{}{
								"name": "test",
							},
							"spec": map[string]interface{}{
								"containers": []interface{}{
									map[string]interface{}{
										"name":  "test-container",
										"image": "test-image",
									},
								},
							},
						},
					},
				},
			},
			expected: 3,
			found:    true,
		},
		{
			name: "extract replicas from not deployment like object",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "SomeTemplatingAPI",
					"spec": map[string]interface{}{
						"replicas": int64(3),
						"selector": map[string]interface{}{
							"matchLabels": map[string]interface{}{
								"app": "test",
							},
						},
						"template": map[string]interface{}{
							"metadata": map[string]interface{}{
								"name": "test",
							},
							"spec": map[string]interface{}{
								"containers": []interface{}{
									map[string]interface{}{
										"name":  "test-container",
										"image": "test-image",
									},
								},
							},
						},
					},
				},
			},
			expected: 0,
			found:    false,
		},
	}

	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			replicas, found := Replicas(tc.obj)
			assert.Equal(t, tc.found, found)
			assert.Equal(t, tc.expected, replicas)
		})
	}
}
