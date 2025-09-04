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
	t.Run("Annotation searched for exists on template, returns no diagnostics", func(t *testing.T) {

		// GIVEN
		// Setup a StatefulSet to check with our new linter
		sts := &appsv1.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{Name: "statefulset"},
			Spec: appsv1.StatefulSetSpec{
				VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
					{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{"valid": "true"}}},
				},
			},
		}

		// SUT Setup
		// Creating mock lint context
		mockLintCtx := mocks.NewMockContext()
		mockLintCtx.AddObject("statefulset", sts)

		// Fetching template
		template, found := templates.Get("statefulset-volumeclaimtemplate-annotation")
		if !found {
			t.Fatalf("failed to get template")
		}

		// Parsing and validating parameters
		params, err := params.ParseAndValidate(map[string]interface{}{"annotation": "valid"})
		if err != nil {
			t.Fatalf("failed to parse and validate params: %v", err)
		}

		// Instantiating check function
		checkFunc, err := template.Instantiate(params)
		if err != nil {
			t.Fatalf("failed to instantiate check function: %v", err)
		}

		// WHEN
		// Running the check function
		diags := checkFunc(mockLintCtx, mockLintCtx.Objects()[0])

		// THEN
		if len(diags) != 0 {
			t.Errorf("got %d diagnostics, want %d", len(diags), 0)
		}
	})
	t.Run("Annotation searched for does not exist on template with annotations, returns a diagnostic", func(t *testing.T) {

		// GIVEN
		// Setup a StatefulSet to check with our new linter
		sts := &appsv1.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{Name: "statefulset"},
			Spec: appsv1.StatefulSetSpec{
				VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
					{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{"annotation": "true"}}},
				},
			},
		}

		// SUT Setup
		// Creating mock lint context
		mockLintCtx := mocks.NewMockContext()
		mockLintCtx.AddObject("statefulset", sts)

		// Fetching template
		template, found := templates.Get("statefulset-volumeclaimtemplate-annotation")
		if !found {
			t.Fatalf("failed to get template")
		}

		// Parsing and validating parameters
		params, err := params.ParseAndValidate(map[string]interface{}{"annotation": "valid"})
		if err != nil {
			t.Fatalf("failed to parse and validate params: %v", err)
		}

		// Instantiating check function
		checkFunc, err := template.Instantiate(params)
		if err != nil {
			t.Fatalf("failed to instantiate check function: %v", err)
		}

		// WHEN
		// Running the check function
		diags := checkFunc(mockLintCtx, mockLintCtx.Objects()[0])

		// THEN
		if len(diags) != 1 {
			t.Errorf("got %d diagnostics, want %d", len(diags), 1)
		}
	})

	t.Run("Annotation searched for does not exist on template without annotations, returns a diagnostic", func(t *testing.T) {

		// GIVEN
		// Setup a StatefulSet to check with our new linter
		sts := &appsv1.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{Name: "statefulset"},
			Spec: appsv1.StatefulSetSpec{
				VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
					{ObjectMeta: metav1.ObjectMeta{}},
				},
			},
		}

		// SUT Setup
		// Creating mock lint context
		mockLintCtx := mocks.NewMockContext()
		mockLintCtx.AddObject("statefulset", sts)

		// Fetching template
		template, found := templates.Get("statefulset-volumeclaimtemplate-annotation")
		if !found {
			t.Fatalf("failed to get template")
		}

		// Parsing and validating parameters
		params, err := params.ParseAndValidate(map[string]interface{}{"annotation": "valid"})
		if err != nil {
			t.Fatalf("failed to parse and validate params: %v", err)
		}

		// Instantiating check function
		checkFunc, err := template.Instantiate(params)
		if err != nil {
			t.Fatalf("failed to instantiate check function: %v", err)
		}

		// WHEN
		// Running the check function
		diags := checkFunc(mockLintCtx, mockLintCtx.Objects()[0])

		// THEN
		if len(diags) != 1 {
			t.Errorf("got %d diagnostics, want %d", len(diags), 1)
		}
	})
}
