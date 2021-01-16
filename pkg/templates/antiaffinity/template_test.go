package antiaffinity

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/internal/pointers"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/antiaffinity/internal/params"
	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAntiAffinity(t *testing.T) {
	suite.Run(t, new(AntiAffinityTestSuite))
}

type AntiAffinityTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *AntiAffinityTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *AntiAffinityTestSuite) addDeploymentWithReplicas(name string, replicas int32) {
	s.ctx.AddMockDeployment(s.T(), name)
	s.ctx.ModifyDeployment(s.T(), name, func(deployment *appsV1.Deployment) {
		deployment.Spec.Replicas = pointers.Int32(replicas)
	})
}

func (s *AntiAffinityTestSuite) TestNoAntiAffinity() {
	const (
		noExplicitReplicasDepName = "no-explicit-replicas"
		oneReplicaDepName         = "one-replica"
		twoReplicasDepName        = "two-replicas"
	)
	s.ctx.AddMockDeployment(s.T(), noExplicitReplicasDepName)
	s.addDeploymentWithReplicas(oneReplicaDepName, 1)
	s.addDeploymentWithReplicas(twoReplicasDepName, 2)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				MinReplicas: 2,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				twoReplicasDepName: {
					{Message: "object has 2 replicas but does not specify inter pod anti-affinity"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				MinReplicas: 1,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				// Replicas defaults to 1.
				noExplicitReplicasDepName: {
					{Message: "object has 1 replica but does not specify inter pod anti-affinity"},
				},
				oneReplicaDepName: {
					{Message: "object has 1 replica but does not specify inter pod anti-affinity"},
				},
				twoReplicasDepName: {
					{Message: "object has 2 replicas but does not specify inter pod anti-affinity"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *AntiAffinityTestSuite) addDeploymentWithAntiAffinity(name string, replicas int32, topologyKey string) {
	s.addDeploymentWithReplicas(name, replicas)
	s.ctx.ModifyDeployment(s.T(), name, func(deployment *appsV1.Deployment) {
		deployment.Spec.Template.Labels = map[string]string{"app": name}
		deployment.Spec.Template.Spec.Affinity = &v1.Affinity{
			PodAntiAffinity: &v1.PodAntiAffinity{
				PreferredDuringSchedulingIgnoredDuringExecution: []v1.WeightedPodAffinityTerm{
					{
						Weight: 1,
						PodAffinityTerm: v1.PodAffinityTerm{
							TopologyKey:   topologyKey,
							LabelSelector: &metaV1.LabelSelector{MatchLabels: map[string]string{"app": name}},
						},
					},
				},
			},
		}
	})

}

func (s *AntiAffinityTestSuite) TestWithAntiAffinity() {
	const (
		kubernetesIOHostnameDepName = "kubernetes-io-hostname"
		otherValidKeyDepName        = "other-valid-key"
		weirdKeyDepName             = "weird-key"
	)
	s.addDeploymentWithAntiAffinity(kubernetesIOHostnameDepName, 2, "kubernetes.io/hostname")
	s.addDeploymentWithAntiAffinity(otherValidKeyDepName, 3, "other.valid/key")
	s.addDeploymentWithAntiAffinity(weirdKeyDepName, 4, "weird/key")

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				MinReplicas: 2,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				otherValidKeyDepName: {
					{Message: "object has 3 replicas but does not specify inter pod anti-affinity"},
				},
				weirdKeyDepName: {
					{Message: "object has 4 replicas but does not specify inter pod anti-affinity"},
				},
			},
			ExpectInstantiationError: false,
		},
		{
			Param: params.Params{
				MinReplicas: 2,
				TopologyKey: "other.valid/key",
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				kubernetesIOHostnameDepName: {
					{Message: "object has 2 replicas but does not specify inter pod anti-affinity"},
				},
				weirdKeyDepName: {
					{Message: "object has 4 replicas but does not specify inter pod anti-affinity"},
				},
			},
			ExpectInstantiationError: false,
		},
		{
			Param: params.Params{
				MinReplicas: 2,
				TopologyKey: ".+",
			},
			Diagnostics:              nil,
			ExpectInstantiationError: false,
		},
	})
}
