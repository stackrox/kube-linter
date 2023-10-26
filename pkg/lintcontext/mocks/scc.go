package mocks

import (
	"testing"

	ocpSecV1 "github.com/openshift/api/security/v1"
	"github.com/stretchr/testify/require"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AddMockSecurityContextConstraints adds a mock SecurityContextConstraints to LintContext
func (l *MockLintContext) AddMockSecurityContextConstraints(t *testing.T, name string, allowFlag bool) {
	require.NotEmpty(t, name)
	l.objects[name] = &ocpSecV1.SecurityContextConstraints{
		TypeMeta: metaV1.TypeMeta{
			Kind:       objectkinds.SecurityContextConstraints,
			APIVersion: objectkinds.GetSCCAPIVersion(),
		},
		ObjectMeta:               metaV1.ObjectMeta{Name: name},
		AllowPrivilegedContainer: allowFlag,
	}
}
