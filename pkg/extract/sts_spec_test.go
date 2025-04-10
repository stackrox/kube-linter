package extract

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.stackrox.io/kube-linter/pkg/k8sutil"
	appsV1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestStatefulSetSpec(t *testing.T) {
	tests := []struct {
		name           string
		obj            k8sutil.Object
		expectedResult appsV1.StatefulSetSpec
		expectedFound  bool
	}{
		{
			name: "StatefulSet object",
			obj: &appsV1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Name: "test-statefulset"},
				Spec: appsV1.StatefulSetSpec{
					Replicas: new(int32),
				},
			},
			expectedResult: appsV1.StatefulSetSpec{
				Replicas: new(int32),
			},
			expectedFound: true,
		},
		{
			name:           "Non-StatefulSet object (Deployment)",
			obj:            &appsV1.Deployment{}, // Geçerli bir objeyi başka bir türde veriyoruz
			expectedResult: appsV1.StatefulSetSpec{},
			expectedFound:  false,
		},
		{
			name:           "Nil object",
			obj:            nil, // Nil obje durumu
			expectedResult: appsV1.StatefulSetSpec{},
			expectedFound:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, found := StatefulSetSpec(tt.obj)
			assert.Equal(t, tt.expectedFound, found)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
