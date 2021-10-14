package replicas

import (
	"testing"

	ocsAppsV1 "github.com/openshift/api/apps/v1"
	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/replicas/internal/params"
	appsV1 "k8s.io/api/apps/v1"
)

func TestReplicas(t *testing.T) {
	suite.Run(t, new(ReplicaTestSuite))
}

type ReplicaTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *ReplicaTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *ReplicaTestSuite) addDeploymentWithReplicas(name string, replicas int32) {
	s.ctx.AddMockDeployment(s.T(), name)
	s.ctx.ModifyDeployment(s.T(), name, func(deployment *appsV1.Deployment) {
		deployment.Spec.Replicas = &replicas
	})
}

func (s *ReplicaTestSuite) addDeploymentConfigWithReplicas(name string, replicas int32) {
	s.ctx.AddMockDeploymentConfig(s.T(), name)
	s.ctx.ModifyDeploymentConfig(s.T(), name, func(dc *ocsAppsV1.DeploymentConfig) {
		dc.Spec.Replicas = replicas
	})
}

func (s *ReplicaTestSuite) TestTooFewReplicas() {
	const (
		noExplicitReplicasDepName = "no-explicit-replicas"
		twoReplicasDepName        = "two-replicas"
		twoReplicasDepNameDc      = "two-replicas-dc"
	)

	s.ctx.AddMockDeployment(s.T(), noExplicitReplicasDepName)
	s.addDeploymentWithReplicas(twoReplicasDepName, 2)

	s.addDeploymentConfigWithReplicas(twoReplicasDepNameDc, 2)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				MinReplicas: 3,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				noExplicitReplicasDepName: {
					{Message: "object has 1 replica but minimum required replicas is 3"},
				},
				twoReplicasDepName: {
					{Message: "object has 2 replicas but minimum required replicas is 3"},
				},
				twoReplicasDepNameDc: {
					{Message: "object has 2 replicas but minimum required replicas is 3"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *ReplicaTestSuite) TestAcceptableReplicas() {
	const (
		acceptableReplicasDepName = "acceptable-replicas"
	)
	s.addDeploymentWithReplicas(acceptableReplicasDepName, 3)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				MinReplicas: 3,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				acceptableReplicasDepName: nil,
			},
			ExpectInstantiationError: false,
		},
	})
}
