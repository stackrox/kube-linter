package pdbminavailable

import (
	"testing"

	kedaV1Alpha1 "github.com/kedacore/keda/v2/apis/keda/v1alpha1"
	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/internal/pointers"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	appsV1 "k8s.io/api/apps/v1"
	autoscalingV2 "k8s.io/api/autoscaling/v2"
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
	tests := []struct {
		name           string
		deploymentSpec appsV1.DeploymentSpec
		pdbSpec        v1.PodDisruptionBudgetSpec
	}{
		{
			name: "replicas equal with matching labels",
			deploymentSpec: appsV1.DeploymentSpec{
				Replicas: pointers.Int32(1),
				Selector: &metaV1.LabelSelector{
					MatchLabels: map[string]string{"foo": "bar"},
				},
			},
			pdbSpec: v1.PodDisruptionBudgetSpec{
				Selector: &metaV1.LabelSelector{
					MatchLabels: map[string]string{"foo": "bar"},
				},
				MinAvailable: &intstr.IntOrString{IntVal: 1},
			},
		},
		{
			name: "replicas equal with matching expression",
			deploymentSpec: appsV1.DeploymentSpec{
				Replicas: pointers.Int32(1),
				Selector: &metaV1.LabelSelector{
					MatchLabels: map[string]string{"foo": "bar"},
				},
			},
			pdbSpec: v1.PodDisruptionBudgetSpec{
				Selector: &metaV1.LabelSelector{
					MatchExpressions: []metaV1.LabelSelectorRequirement{
						{
							Key:      "foo",
							Operator: metaV1.LabelSelectorOpIn,
							Values:   []string{"baz", "bar", "qux"},
						},
					},
				},
				MinAvailable: &intstr.IntOrString{IntVal: 1},
			},
		},
	}

	for _, tt := range tests {
		p.T().Run(tt.name, func(t *testing.T) {
			p.ctx.AddMockDeployment(p.T(), "test-deploy")
			p.ctx.ModifyDeployment(p.T(), "test-deploy", func(deployment *appsV1.Deployment) {
				deployment.Namespace = "test"
				deployment.Spec = tt.deploymentSpec
			})
			p.ctx.AddMockPodDisruptionBudget(p.T(), "test-pdb")
			p.ctx.ModifyPodDisruptionBudget(p.T(), "test-pdb", func(pdb *v1.PodDisruptionBudget) {
				pdb.Namespace = "test"
				pdb.Spec = tt.pdbSpec
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
		})
	}
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

// test that the check run with a deployment that has no replicas and a HPA that has a minReplicas
func (p *PDBTestSuite) TestPDBWithMinAvailableHPA() {
	p.ctx.AddMockDeployment(p.T(), "test-deploy")
	p.ctx.ModifyDeployment(p.T(), "test-deploy", func(deployment *appsV1.Deployment) {
		deployment.Namespace = "test"
		deployment.Spec.Replicas = nil
		deployment.Spec.Selector = &metaV1.LabelSelector{}
		deployment.Spec.Selector.MatchLabels = map[string]string{"foo": "bar"}
	})
	p.ctx.AddMockHorizontalPodAutoscaler(p.T(), "test-hpa", "v2")
	p.ctx.ModifyHorizontalPodAutoscalerV2(p.T(), "test-hpa", func(hpa *autoscalingV2.HorizontalPodAutoscaler) {
		hpa.Namespace = "test"
		hpa.Spec.ScaleTargetRef = autoscalingV2.CrossVersionObjectReference{
			Kind:       "Deployment",
			Name:       "test-deploy",
			APIVersion: "apps/v1",
		}
		hpa.Spec.MinReplicas = nil
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
				"test-pdb": {
					{Message: "The current number of replicas for deployment test-deploy is equal to or lower than the minimum number of replicas specified by its PDB."},
				},
			},
			ExpectInstantiationError: false,
		},
	})

}

// test that the check run with a deployment that has no replicas and a Keda ScaledObject that don't has a minReplicas
func (p *PDBTestSuite) TestPDBWithMinAvailableAndKedaScaledObjectDoNotHasMinReplicas() {
	p.ctx.AddMockDeployment(p.T(), "test-deploy")
	p.ctx.ModifyDeployment(p.T(), "test-deploy", func(deployment *appsV1.Deployment) {
		deployment.Namespace = "test"
		deployment.Spec.Replicas = nil
		deployment.Spec.Selector = &metaV1.LabelSelector{}
		deployment.Spec.Selector.MatchLabels = map[string]string{"foo": "bar"}
	})
	p.ctx.AddMockScaledObject(p.T(), "test-scaledobject", "v1alpha1")
	p.ctx.ModifyScaledObjectV1Alpha1(p.T(), "test-scaledobject", func(scaledobject *kedaV1Alpha1.ScaledObject) {
		scaledobject.Namespace = "test"
		scaledobject.Spec.ScaleTargetRef = &kedaV1Alpha1.ScaleTarget{
			Kind:       "Deployment",
			Name:       "test-deploy",
			APIVersion: "apps/v1",
		}
		scaledobject.Spec.MinReplicaCount = nil
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
				"test-pdb": {
					{Message: "The current number of replicas for deployment test-deploy is equal to or lower than the minimum number of replicas specified by its PDB."},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

// test that the check run with a deployment that has no replicas and a Keda ScaledObject that has a minReplicas
func (p *PDBTestSuite) TestPDBWithMinAvailableAndKedaScaledObjectHasMinReplicas() {
	p.ctx.AddMockDeployment(p.T(), "test-deploy")
	p.ctx.ModifyDeployment(p.T(), "test-deploy", func(deployment *appsV1.Deployment) {
		deployment.Namespace = "test"
		deployment.Spec.Replicas = nil
		deployment.Spec.Selector = &metaV1.LabelSelector{}
		deployment.Spec.Selector.MatchLabels = map[string]string{"foo": "bar"}
	})
	p.ctx.AddMockScaledObject(p.T(), "test-scaledobject", "v1alpha1")
	p.ctx.ModifyScaledObjectV1Alpha1(p.T(), "test-scaledobject", func(scaledobject *kedaV1Alpha1.ScaledObject) {
		scaledobject.Namespace = "test"
		scaledobject.Spec.ScaleTargetRef = &kedaV1Alpha1.ScaleTarget{
			Kind:       "Deployment",
			Name:       "test-deploy",
			APIVersion: "apps/v1",
		}
		scaledobject.Spec.MinReplicaCount = pointers.Int32(4)
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
