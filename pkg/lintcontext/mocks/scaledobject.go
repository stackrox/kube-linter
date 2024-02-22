package mocks

import (
	"fmt"
	"testing"

	kedaV1Alpha1 "github.com/kedacore/keda/v2/apis/keda/v1alpha1"
	"github.com/stretchr/testify/require"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AddMockScaledObject adds a mock ScaledObject to LintContext
func (l *MockLintContext) AddMockScaledObject(t *testing.T, name, version string) {
	require.NotEmpty(t, name)
	switch version {
	case "v1alpha1":
		l.objects[name] = &kedaV1Alpha1.ScaledObject{
			TypeMeta: metaV1.TypeMeta{
				Kind:       objectkinds.ScaledObject,
				APIVersion: objectkinds.GetScaledObjectAPIVersion(version),
			},
			ObjectMeta: metaV1.ObjectMeta{Name: name},
			Spec:       kedaV1Alpha1.ScaledObjectSpec{},
		}
	default:
		require.FailNow(t, fmt.Sprintf("Unknown scaled object version %s", version))
	}
}

// ModifyScaledObjectV1Alpha1 modifies a given ScaledObject in the context via the passed function.
func (l *MockLintContext) ModifyScaledObjectV1Alpha1(t *testing.T, name string, f func(hpa *kedaV1Alpha1.ScaledObject)) {
	r, ok := l.objects[name].(*kedaV1Alpha1.ScaledObject)
	require.True(t, ok)
	f(r)
}
