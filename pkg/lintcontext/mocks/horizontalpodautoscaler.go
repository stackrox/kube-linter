package mocks

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	autoscalingV2Beta1 "k8s.io/api/autoscaling/v2beta1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AddMockHorizontalPodAutoscaler adds a mock HorizontalPodAutoscaler to LintContext
func (l *MockLintContext) AddMockHorizontalPodAutoscaler(t *testing.T, name string) {
	require.NotEmpty(t, name)
	l.objects[name] = &autoscalingV2Beta1.HorizontalPodAutoscaler{
		TypeMeta: metaV1.TypeMeta{
			Kind:       objectkinds.HorizontalPodAutoscaler,
			APIVersion: objectkinds.GetHorizontalPodAutoscalerAPIVersion(),
		},
		ObjectMeta: metaV1.ObjectMeta{Name: name},
		Spec:       autoscalingV2Beta1.HorizontalPodAutoscalerSpec{},
	}
}

//ModifyHorizontalPodAutoscaler modifies a given HorizontalPodAutoscaler in the context via the passed function.
func (l *MockLintContext) ModifyHorizontalPodAutoscaler(t *testing.T, name string, f func(hpa *autoscalingV2Beta1.HorizontalPodAutoscaler)) {
	r, ok := l.objects[name].(*autoscalingV2Beta1.HorizontalPodAutoscaler)
	require.True(t, ok)
	f(r)
}
