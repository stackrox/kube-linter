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

func (s *AntiAffinityTestSuite) TestEmptyAntiAffinity() {
	const (
		oneReplicaDepName  = "one-replica"
		twoReplicasDepName = "two-replicas"
	)

	s.addDeploymentWithEmptyAntiAffinity(oneReplicaDepName, 1)
	s.addDeploymentWithEmptyAntiAffinity(twoReplicasDepName, 2)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				MinReplicas: 1,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				oneReplicaDepName: {
					{Message: "object has 1 replica but does not specify preferred or required inter " +
						"pod anti-affinity during scheduling"},
				},
				twoReplicasDepName: {
					{Message: "object has 2 replicas but does not specify preferred or required inter " +
						"pod anti-affinity during scheduling"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *AntiAffinityTestSuite) addDeploymentWithAntiAffinity(name string, replicas int32, topologyKey string,
	labelName string, namespace string) {
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
							Namespaces:    []string{namespace},
							LabelSelector: &metaV1.LabelSelector{MatchLabels: map[string]string{"app": labelName}},
						},
					},
				},
			},
		}
	})
}

func (s *AntiAffinityTestSuite) addDeploymentWithEmptyAntiAffinity(name string, replicas int32) {
	s.addDeploymentWithReplicas(name, replicas)
	s.ctx.ModifyDeployment(s.T(), name, func(deployment *appsV1.Deployment) {
		deployment.Spec.Template.Spec.Affinity = &v1.Affinity{
			PodAntiAffinity: &v1.PodAntiAffinity{
				PreferredDuringSchedulingIgnoredDuringExecution: []v1.WeightedPodAffinityTerm{},
				RequiredDuringSchedulingIgnoredDuringExecution:  []v1.PodAffinityTerm{},
			},
		}
	})
}

func (s *AntiAffinityTestSuite) TestWithAntiAffinity() {
	const (
		kubernetesIOHostnameDepName = "kubernetes-io-hostname"
		otherValidKeyDepName        = "other-valid-key"
		weirdKeyDepName             = "weird-key"
		nonMatchingLabelSelectors   = "non-matching-label-selector"
		nonMatchingNamespace        = "non-matching-namespace"
	)
	s.addDeploymentWithAntiAffinity(kubernetesIOHostnameDepName, 2, "kubernetes.io/hostname",
		kubernetesIOHostnameDepName, "")
	s.addDeploymentWithAntiAffinity(otherValidKeyDepName, 3, "other.valid/key",
		otherValidKeyDepName, "")
	s.addDeploymentWithAntiAffinity(weirdKeyDepName, 4, "weird/key", weirdKeyDepName, "")
	s.addDeploymentWithAntiAffinity(nonMatchingLabelSelectors, 4, "other.valid/key",
		"non-matching", "")
	s.addDeploymentWithAntiAffinity(nonMatchingNamespace, 4, "other.valid/key",
		nonMatchingNamespace, "non-matching-namespace")

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				MinReplicas: 2,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				otherValidKeyDepName: {
					{Message: "anti-affinity's topology key does not match \"other.valid/key\""},
				},
				weirdKeyDepName: {
					{Message: "anti-affinity's topology key does not match \"weird/key\""},
				},
				nonMatchingLabelSelectors: {
					{Message: "anti-affinity's topology key does not match \"other.valid/key\""},
				},
				nonMatchingNamespace: {
					{Message: "pod's namespace \"\" not found in anti-affinity's namespaces [non-matching-namespace]"},
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
					{Message: "anti-affinity's topology key does not match \"kubernetes.io/hostname\""},
				},
				weirdKeyDepName: {
					{Message: "anti-affinity's topology key does not match \"weird/key\""},
				},
				nonMatchingLabelSelectors: {
					{Message: "pod's labels \"app=non-matching-label-selector\" do not match with anti-affinity's " +
						"labels \"app=non-matching\""},
				},
				nonMatchingNamespace: {
					{Message: "pod's namespace \"\" not found in anti-affinity's namespaces [non-matching-namespace]"},
				},
			},
			ExpectInstantiationError: false,
		},
		{
			Param: params.Params{
				MinReplicas: 2,
				TopologyKey: ".+",
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				nonMatchingLabelSelectors: {
					{Message: "pod's labels \"app=non-matching-label-selector\" do not match with anti-affinity's " +
						"labels \"app=non-matching\""},
				},
				nonMatchingNamespace: {
					{Message: "pod's namespace \"\" not found in anti-affinity's namespaces [non-matching-namespace]"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}
