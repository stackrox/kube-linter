package extract

import (
	"testing"
	"time"

	appsV1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/stretchr/testify/assert"
)

type fakeStatefulSet struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec appsV1.StatefulSetSpec
}

func (f *fakeStatefulSet) GetObjectKind() schema.ObjectKind {
	return &f.TypeMeta
}

func (f *fakeStatefulSet) DeepCopyObject() runtime.Object {
	return &fakeStatefulSet{
		TypeMeta: f.TypeMeta,
		ObjectMeta: metav1.ObjectMeta{
			Name:      f.Name,
			Namespace: f.Namespace,
		},
		Spec: f.Spec,
	}
}

func (f *fakeStatefulSet) GetAnnotations() map[string]string {
	return map[string]string{"key": "value"} // Example annotation
}

func (f *fakeStatefulSet) GetCreationTimestamp() metav1.Time {
	return metav1.Time{Time: time.Now()}
}

func TestStatefulSetSpec(t *testing.T) {
	t.Run("nil object", func(t *testing.T) {
		spec, ok := StatefulSetSpec(nil)
		assert.False(t, ok)
		assert.Equal(t, appsV1.StatefulSetSpec{}, spec)
	})

	t.Run("typed StatefulSet", func(t *testing.T) {
		sampleSpec := appsV1.StatefulSetSpec{
			ServiceName: "my-service",
		}
		obj := &appsV1.StatefulSet{
			Spec: sampleSpec,
		}
		spec, ok := StatefulSetSpec(obj)
		assert.True(t, ok)
		assert.Equal(t, sampleSpec, spec)
	})

	t.Run("fallback via reflection", func(t *testing.T) {
		sampleSpec := appsV1.StatefulSetSpec{
			ServiceName: "reflected-service",
		}
		obj := &fakeStatefulSet{
			TypeMeta: metav1.TypeMeta{
				Kind: "StatefulSet",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "fake-statefulset",
				Namespace: "default",
			},
			Spec: sampleSpec,
		}
		spec, ok := StatefulSetSpec(obj)
		assert.True(t, ok)
		assert.Equal(t, sampleSpec, spec)
	})

	t.Run("wrong kind", func(t *testing.T) {
		obj := &fakeStatefulSet{
			TypeMeta: metav1.TypeMeta{
				Kind: "Deployment",
			},
			Spec: appsV1.StatefulSetSpec{},
		}
		spec, ok := StatefulSetSpec(obj)
		assert.False(t, ok)
		assert.Equal(t, appsV1.StatefulSetSpec{}, spec)
	})
}
