package volumeclaimtemplates

import (
	"testing"

	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/volumeclaimtemplates/internal/params"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestStatefulSetVolumeClaimTemplateAnnotation(t *testing.T) {
	tests := []struct {
		name        string
		annotation  string  // Adjusted to match the parameter name in Params
		wantDiags   int
	}{
		{"WithAnnotation", "value", 0},
		{"WithoutAnnotation", "", 1},  // Empty string or any value for negative case
	}
	
	

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sts := &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Name: "statefulset"},
				Spec: appsv1.StatefulSetSpec{
					VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
						{ObjectMeta: metav1.ObjectMeta{Annotations: tt.annotations}},
					},
				},
			}

			// Creating mock lint context
			mockLintCtx := mocks.NewMockContext()
			mockLintCtx.AddObject("statefulset", sts)

			// Fetching template
			template, found := templates.Get("statefulset-volumeclaimtemplate-annotation")
			if !found {
				t.Fatalf("failed to get template")
			}

			// Parsing and validating parameters
			params, err := params.ParseAndValidate(map[string]interface{}{})
			if err != nil {
				t.Fatalf("failed to parse and validate params: %v", err)
			}

			// Instantiating check function
			checkFunc, err := template.Instantiate(params)
			if err != nil {
				t.Fatalf("failed to instantiate check function: %v", err)
			}

			// Running the check function
			diags := checkFunc(mockLintCtx, mockLintCtx.Objects()[0])
			if len(diags) != tt.wantDiags {
				t.Errorf("got %d diagnostics, want %d", len(diags), tt.wantDiags)
			}
		})
	}
}
