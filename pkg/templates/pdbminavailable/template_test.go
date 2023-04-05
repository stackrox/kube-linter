package pdbminavailable

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/internal/pointers"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/policy/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"golang.stackrox.io/kube-linter/pkg/templates/pdbminavailable/internal/params"
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

func (p *PDBTestSuite) TestPDBMinAvailableZero() {

	p.ctx.AddMockDeployment(p.T(), "test-deploy")
	p.ctx.ModifyDeployment(p.T(), "test-deploy", func(deployment *appsV1.Deployment) {
		deployment.Namespace = "test"
		deployment.Spec.Replicas = pointers.Int32(1)
		deployment.Spec.Selector = &metaV1.LabelSelector{}
		deployment.Spec.Selector.MatchLabels = map[string]string{"foo": "bar"}
	})

	p.ctx.AddMockPodDisruptionBudget(p.T(), "test-pdb")
	p.ctx.ModifyPodDisruptionBudget(p.T(), "test-pdb", func(pdb *v1.PodDisruptionBudget) {
		pdb.Namespace = "test"
		pdb.Spec.Selector = &metaV1.LabelSelector{}
		pdb.Spec.Selector.MatchLabels = map[string]string{"foo": "bar"}
		pdb.Spec.MinAvailable = &intstr.IntOrString{IntVal: 0}
	})

	p.Validate(p.ctx, []templates.TestCase{
		{
			Param:                    params.Params{},
			Diagnostics:              map[string][]diagnostic.Diagnostic{},
			ExpectInstantiationError: false,
		},
	})
}

func (p *PDBTestSuite) TestPDBMinAvailableReplicasEqual() {

	p.ctx.AddMockDeployment(p.T(), "test-deploy")
	p.ctx.ModifyDeployment(p.T(), "test-deploy", func(deployment *appsV1.Deployment) {
		deployment.Namespace = "test"
		deployment.Spec.Replicas = pointers.Int32(1)
		deployment.Spec.Selector = &metaV1.LabelSelector{}
		deployment.Spec.Selector.MatchLabels = map[string]string{"foo": "bar"}
	})

	p.ctx.AddMockPodDisruptionBudget(p.T(), "test-pdb")
	p.ctx.ModifyPodDisruptionBudget(p.T(), "test-pdb", func(pdb *v1.PodDisruptionBudget) {
		pdb.Namespace = "test"
		pdb.Spec.Selector = &metaV1.LabelSelector{}
		pdb.Spec.Selector.MatchLabels = map[string]string{"foo": "bar"}
		pdb.Spec.MinAvailable = &intstr.IntOrString{IntVal: 1}
	})

	p.Validate(p.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				"test-pdb": {
					{Message: "The current number of replicas for deployment test-deploy is equal to or lower than the minimum number of replicas specified by its PDB."},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (p *PDBTestSuite) TestPDBMinAvailableFiftyPercent() {

	p.ctx.AddMockDeployment(p.T(), "test-deploy")
	p.ctx.ModifyDeployment(p.T(), "test-deploy", func(deployment *appsV1.Deployment) {
		deployment.Namespace = "test"
		deployment.Spec.Replicas = pointers.Int32(2)
		deployment.Spec.Selector = &metaV1.LabelSelector{}
		deployment.Spec.Selector.MatchLabels = map[string]string{"foo": "bar"}
	})

	p.ctx.AddMockPodDisruptionBudget(p.T(), "test-pdb")
	p.ctx.ModifyPodDisruptionBudget(p.T(), "test-pdb", func(pdb *v1.PodDisruptionBudget) {
		pdb.Namespace = "test"
		pdb.Spec.Selector = &metaV1.LabelSelector{}
		pdb.Spec.Selector.MatchLabels = map[string]string{"foo": "bar"}
		pdb.Spec.MinAvailable = &intstr.IntOrString{StrVal: "50%", Type: intstr.String}
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

func (p *PDBTestSuite) TestPDBMinAvailableOneHundredPercent() {
	p.ctx.AddMockPodDisruptionBudget(p.T(), "test-pdb")
	p.ctx.ModifyPodDisruptionBudget(p.T(), "test-pdb", func(pdb *v1.PodDisruptionBudget) {
		pdb.Spec.MinAvailable = &intstr.IntOrString{StrVal: "100%", Type: intstr.String}
	})

	p.Validate(p.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				"test-pdb": {
					{Message: "PDB has minimum available replicas set to 100 percent of replicas"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}
