package mocks

import (
	"testing"

	"github.com/stretchr/testify/require"
	pdbV1 "k8s.io/api/policy/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AddMockPodDisruptionBudget adds a mock PodDisruptionBudget to LintContext
func (l *MockLintContext) AddMockPodDisruptionBudget(t *testing.T, name string) {
	require.NotEmpty(t, name)
	l.objects[name] = &pdbV1.PodDisruptionBudget{
		ObjectMeta: metaV1.ObjectMeta{Name: name},
	}
}

// ModifyPodDisruptionBudget modifies a given PodDisruptionBudget in the context via the passed function.
func (l *MockLintContext) ModifyPodDisruptionBudget(t *testing.T, name string, f func(podDisruptionBudget *pdbV1.PodDisruptionBudget)) {
	pdb, ok := l.objects[name].(*pdbV1.PodDisruptionBudget)
	require.True(t, ok)
	f(pdb)
}
