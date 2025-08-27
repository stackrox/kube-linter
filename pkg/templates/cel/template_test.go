package cel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name        string
		check       string
		subject     lintcontext.Object
		objects     []lintcontext.Object
		expectedMsg string
		expectError bool
	}{
		{
			name:  "simple string return",
			check: `"test message"`,
			subject: lintcontext.Object{
				K8sObject: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
				},
			},
			objects:     []lintcontext.Object{},
			expectedMsg: "test message",
			expectError: false,
		},
		{
			name:  "conditional message based on subject",
			check: `subject.metadata.name == "test-pod" ? "pod found" : "pod not found"`,
			subject: lintcontext.Object{
				K8sObject: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
				},
			},
			objects:     []lintcontext.Object{},
			expectedMsg: "pod found",
			expectError: false,
		},
		{
			name:  "check objects list length",
			check: `size(objects) > 0 ? "objects found" : "no objects"`,
			subject: lintcontext.Object{
				K8sObject: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
				},
			},
			objects: []lintcontext.Object{
				{
					K8sObject: &corev1.Service{
						ObjectMeta: metav1.ObjectMeta{Name: "test-service"},
					},
				},
			},
			expectedMsg: "objects found",
			expectError: false,
		},
		{
			name:  "complex expression with subject properties",
			check: `subject.metadata.namespace == "default" && subject.metadata.name.startsWith("test") ? "valid test object" : "invalid object"`,
			subject: lintcontext.Object{
				K8sObject: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-pod",
						Namespace: "default",
					},
				},
			},
			objects:     []lintcontext.Object{},
			expectedMsg: "valid test object",
			expectError: false,
		},
		{
			name:  "empty string return",
			check: `""`,
			subject: lintcontext.Object{
				K8sObject: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
				},
			},
			objects:     []lintcontext.Object{},
			expectedMsg: "",
			expectError: false,
		},
		{
			name:  "invalid CEL expression",
			check: `invalid syntax here`,
			subject: lintcontext.Object{
				K8sObject: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
				},
			},
			objects:     []lintcontext.Object{},
			expectedMsg: "",
			expectError: true,
		},
		{
			name:  "non-string return type",
			check: `123`,
			subject: lintcontext.Object{
				K8sObject: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
				},
			},
			objects:     []lintcontext.Object{},
			expectedMsg: "",
			expectError: true,
		},
		{
			name:  "boolean return type",
			check: `true`,
			subject: lintcontext.Object{
				K8sObject: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
				},
			},
			objects:     []lintcontext.Object{},
			expectedMsg: "",
			expectError: true,
		},
		{
			name:  "accessing nested properties",
			check: `subject.metadata.labels.app == "web" ? "web app detected" : "not a web app"`,
			subject: lintcontext.Object{
				K8sObject: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-pod",
						Labels: map[string]string{
							"app": "web",
						},
					},
				},
			},
			objects:     []lintcontext.Object{},
			expectedMsg: "web app detected",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := evaluate(tt.check, tt.subject, tt.objects)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, msg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMsg, msg)
			}
		})
	}
}
