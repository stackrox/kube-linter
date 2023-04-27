package pdbmaxunavailable

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	pdbv1 "k8s.io/api/policy/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"golang.stackrox.io/kube-linter/pkg/templates/pdbmaxunavailable/internal/params"
)

func TestPDBs(t *testing.T) {
	suite.Run(t, new(PDBTestSuite))
}

type PDBTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (p *PDBTestSuite) SetupTest() {
	p.Init(templateKey)
	p.ctx = mocks.NewMockContext()
}

func toPointer(v intstr.IntOrString) *intstr.IntOrString {
	return &v
}

func (p *PDBTestSuite) TestPDB() {
	testCases := []struct {
		Name           string
		Message        string
		MaxUnavailable *intstr.IntOrString
	}{
		{
			Name:           "As Not Defined (nil)",
			MaxUnavailable: nil,
		},
		{
			Name:           "Invalid Value As String",
			Message:        "maxUnavailable has invalid value [invalid-value]",
			MaxUnavailable: toPointer(intstr.FromString("invalid-value")),
		},
		{
			Name:           "Zero As Integer",
			Message:        "MaxUnavailable is set to 0",
			MaxUnavailable: toPointer(intstr.FromInt(0)),
		},
		{
			Name:           "Zero As String Percentage",
			Message:        "MaxUnavailable is set to 0",
			MaxUnavailable: toPointer(intstr.FromString("0%")),
		},
		{
			Name:           "Zero As String",
			Message:        "maxUnavailable has invalid value [0]",
			MaxUnavailable: toPointer(intstr.FromString("0")),
		},
		{
			Name:           "One As Integer",
			MaxUnavailable: toPointer(intstr.FromInt(1)),
		},
		{
			Name:           "One As String Percentage",
			MaxUnavailable: toPointer(intstr.FromString("1%")),
		},
	}

	for _, tc := range testCases {
		tc := tc
		p.Run(tc.Name, func() {
			p.ctx.AddMockPodDisruptionBudget(p.T(), "test-pdb")
			p.ctx.ModifyPodDisruptionBudget(p.T(), "test-pdb", func(pdb *pdbv1.PodDisruptionBudget) {
				pdb.Spec.MaxUnavailable = tc.MaxUnavailable
			})

			expected := map[string][]diagnostic.Diagnostic{}
			if tc.Message != "" {
				expected = map[string][]diagnostic.Diagnostic{
					"test-pdb": {{Message: tc.Message}},
				}
			}

			p.Validate(p.ctx, []templates.TestCase{
				{
					Param:                    params.Params{},
					Diagnostics:              expected,
					ExpectInstantiationError: false,
				},
			})
		})
	}
}
