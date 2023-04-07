package pdbmaxunavailable

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	v1 "k8s.io/api/policy/v1"
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

func (p *PDBTestSuite) TestPDBMaxUnavailableZero() {

	p.ctx.AddMockPodDisruptionBudget(p.T(), "test-pdb")
	p.ctx.ModifyPodDisruptionBudget(p.T(), "test-pdb", func(pdb *v1.PodDisruptionBudget) {
		pdb.Spec.MaxUnavailable = &intstr.IntOrString{IntVal: 0}
	})

	p.Validate(p.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				"test-pdb": {
					{Message: "MaxUnavailable is set to 0"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (p *PDBTestSuite) TestPDBMaxUnavailableOne() {

	p.ctx.AddMockPodDisruptionBudget(p.T(), "test-pdb")
	p.ctx.ModifyPodDisruptionBudget(p.T(), "test-pdb", func(pdb *v1.PodDisruptionBudget) {
		pdb.Spec.MaxUnavailable = &intstr.IntOrString{IntVal: 1}
	})

	p.Validate(p.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				"test-pdb": {},
			},
			ExpectInstantiationError: false,
		},
	})
}
