package mocks

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	autoscalingV1 "k8s.io/api/autoscaling/v1"
	autoscalingV2 "k8s.io/api/autoscaling/v2"
	autoscalingV2Beta1 "k8s.io/api/autoscaling/v2beta1"
	autoscalingV2Beta2 "k8s.io/api/autoscaling/v2beta2"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AddMockHorizontalPodAutoscaler adds a mock HorizontalPodAutoscaler to LintContext
func (l *MockLintContext) AddMockHorizontalPodAutoscaler(t *testing.T, name, version string) {
	require.NotEmpty(t, name)
	switch version {
	case "v2beta1":
		l.objects[name] = &autoscalingV2Beta1.HorizontalPodAutoscaler{
			TypeMeta: metaV1.TypeMeta{
				Kind:       objectkinds.HorizontalPodAutoscaler,
				APIVersion: objectkinds.GetHorizontalPodAutoscalerAPIVersion(version),
			},
			ObjectMeta: metaV1.ObjectMeta{Name: name},
			Spec:       autoscalingV2Beta1.HorizontalPodAutoscalerSpec{},
		}
	case "v2beta2":
		l.objects[name] = &autoscalingV2Beta2.HorizontalPodAutoscaler{
			TypeMeta: metaV1.TypeMeta{
				Kind:       objectkinds.HorizontalPodAutoscaler,
				APIVersion: objectkinds.GetHorizontalPodAutoscalerAPIVersion(version),
			},
			ObjectMeta: metaV1.ObjectMeta{Name: name},
			Spec:       autoscalingV2Beta2.HorizontalPodAutoscalerSpec{},
		}
	case "v2":
		l.objects[name] = &autoscalingV2.HorizontalPodAutoscaler{
			TypeMeta: metaV1.TypeMeta{
				Kind:       objectkinds.HorizontalPodAutoscaler,
				APIVersion: objectkinds.GetHorizontalPodAutoscalerAPIVersion(version),
			},
			ObjectMeta: metaV1.ObjectMeta{Name: name},
			Spec:       autoscalingV2.HorizontalPodAutoscalerSpec{},
		}
	case "v1":
		l.objects[name] = &autoscalingV1.HorizontalPodAutoscaler{
			TypeMeta: metaV1.TypeMeta{
				Kind:       objectkinds.HorizontalPodAutoscaler,
				APIVersion: objectkinds.GetHorizontalPodAutoscalerAPIVersion(version),
			},
			ObjectMeta: metaV1.ObjectMeta{Name: name},
			Spec:       autoscalingV1.HorizontalPodAutoscalerSpec{},
		}
	default:
		require.FailNow(t, fmt.Sprintf("Unknown autoscaling version %s", version))
	}
}

// ModifyHorizontalPodAutoscalerV2Beta1 modifies a given HorizontalPodAutoscaler in the context via the passed function.
func (l *MockLintContext) ModifyHorizontalPodAutoscalerV2Beta1(t *testing.T, name string, f func(hpa *autoscalingV2Beta1.HorizontalPodAutoscaler)) {
	r, ok := l.objects[name].(*autoscalingV2Beta1.HorizontalPodAutoscaler)
	require.True(t, ok)
	f(r)
}

// ModifyHorizontalPodAutoscalerV2Beta2 modifies a given HorizontalPodAutoscaler in the context via the passed function.
func (l *MockLintContext) ModifyHorizontalPodAutoscalerV2Beta2(t *testing.T, name string, f func(hpa *autoscalingV2Beta2.HorizontalPodAutoscaler)) {
	r, ok := l.objects[name].(*autoscalingV2Beta2.HorizontalPodAutoscaler)
	require.True(t, ok)
	f(r)
}

// ModifyHorizontalPodAutoscalerV2 modifies a given HorizontalPodAutoscaler in the context via the passed function.
func (l *MockLintContext) ModifyHorizontalPodAutoscalerV2(t *testing.T, name string, f func(hpa *autoscalingV2.HorizontalPodAutoscaler)) {
	r, ok := l.objects[name].(*autoscalingV2.HorizontalPodAutoscaler)
	require.True(t, ok)
	f(r)
}

// ModifyHorizontalPodAutoscalerV1 modifies a given HorizontalPodAutoscaler in the context via the passed function.
func (l *MockLintContext) ModifyHorizontalPodAutoscalerV1(t *testing.T, name string, f func(hpa *autoscalingV1.HorizontalPodAutoscaler)) {
	r, ok := l.objects[name].(*autoscalingV1.HorizontalPodAutoscaler)
	require.True(t, ok)
	f(r)
}
