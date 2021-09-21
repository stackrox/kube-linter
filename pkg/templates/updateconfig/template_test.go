package updateconfig

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"golang.stackrox.io/kube-linter/internal/pointers"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/updateconfig/internal/params"

	ocsAppsv1 "github.com/openshift/api/apps/v1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestUpgradeConfig(t *testing.T) {
	suite.Run(t, new(UpgradeConfigTestSuite))
}

type UpgradeConfigTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *UpgradeConfigTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *UpgradeConfigTestSuite) addDeploymentWithStrategy(name string, strategy appsv1.DeploymentStrategy) {
	s.ctx.AddMockDeployment(s.T(), name)
	s.ctx.ModifyDeployment(s.T(), name, func(deployment *appsv1.Deployment) {
		deployment.Spec.Strategy = strategy
	})
}

func (s *UpgradeConfigTestSuite) addDaemonSetWithStrategy(name string, strategy appsv1.DaemonSetUpdateStrategy) {
	s.ctx.AddMockDaemonSet(s.T(), name)
	s.ctx.ModifyDaemonSet(s.T(), name, func(ds *appsv1.DaemonSet) {
		ds.Spec.UpdateStrategy = strategy
	})
}

func (s *UpgradeConfigTestSuite) addDeploymentConfigWithStrategy(name string, strategy ocsAppsv1.DeploymentStrategy) {
	s.ctx.AddMockDeploymentConfig(s.T(), name)
	s.ctx.ModifyDeploymentConfig(s.T(), name, func(ds *ocsAppsv1.DeploymentConfig) {
		ds.Spec.Strategy = strategy
	})
}

func (s *UpgradeConfigTestSuite) addReplicationControllerWithReplicas(name string, replicas int32) {
	s.ctx.AddMockReplicationController(s.T(), name)
	s.ctx.ModifyReplicationController(s.T(), name, func(rc *v1.ReplicationController) {
		rc.Spec.Replicas = pointers.Int32(replicas)
	})
}

func (s *UpgradeConfigTestSuite) TestInvalidStrategyType() {
	const (
		noExplicitStrategy           = "no-explicit-strategy"
		deploymentWithStrategy       = "deployment-strategy-recreate"
		daemonSetWithStrategy        = "daemon-set-strategy-on-delete"
		deploymentConfigWithStrategy = "deployment-config-strategy-recreate"
		replicationController        = "replication-controller"
		deploymentStrategyType       = appsv1.RecreateDeploymentStrategyType
		daemonSetStrategyType        = appsv1.OnDeleteDaemonSetStrategyType
		deploymentConfigStrategyType = ocsAppsv1.DeploymentStrategyTypeRecreate
		strategyRegex                = "^(RollingUpdate|Rolling)$"
	)
	s.ctx.AddMockDeployment(s.T(), noExplicitStrategy)
	s.addDeploymentWithStrategy(deploymentWithStrategy, appsv1.DeploymentStrategy{Type: deploymentStrategyType})
	s.addDaemonSetWithStrategy(daemonSetWithStrategy, appsv1.DaemonSetUpdateStrategy{Type: daemonSetStrategyType})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithStrategy, ocsAppsv1.DeploymentStrategy{Type: deploymentConfigStrategyType})
	s.addReplicationControllerWithReplicas(replicationController, 2)

	deploymentErrorMsg := fmt.Sprintf("object has %s strategy type but must match regex %s", deploymentStrategyType, strategyRegex)
	daemonSetErrorMsg := fmt.Sprintf("object has %s strategy type but must match regex %s", daemonSetStrategyType, strategyRegex)
	deploymentConfigErrorMsg := fmt.Sprintf("object has %s strategy type but must match regex %s", deploymentConfigStrategyType, strategyRegex)
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				StrategyTypeRegex: strategyRegex,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				noExplicitStrategy: {
					{Message: fmt.Sprintf("object has no strategy type but must match regex %s", strategyRegex)},
				},
				deploymentWithStrategy: {
					{Message: deploymentErrorMsg},
				},
				daemonSetWithStrategy: {
					{Message: daemonSetErrorMsg},
				},
				deploymentConfigWithStrategy: {
					{Message: deploymentConfigErrorMsg},
				},
				replicationController: {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *UpgradeConfigTestSuite) TestValidStrategyType() {
	const (
		deploymentWithStrategy       = "deployment-strategy-recreate"
		daemonSetWithStrategy        = "daemon-set-strategy-on-delete"
		deploymentConfigWithStrategy = "deployment-config-strategy-recreate"
		replicationController        = "replication-controller"
		deploymentStrategyType       = appsv1.RollingUpdateDeploymentStrategyType
		daemonSetStrategyType        = appsv1.RollingUpdateDaemonSetStrategyType
		deploymentConfigStrategyType = ocsAppsv1.DeploymentStrategyTypeRolling
		strategyRegex                = "^(RollingUpdate|Rolling)$"
	)
	s.addDeploymentWithStrategy(deploymentWithStrategy, appsv1.DeploymentStrategy{Type: deploymentStrategyType})
	s.addDaemonSetWithStrategy(daemonSetWithStrategy, appsv1.DaemonSetUpdateStrategy{Type: daemonSetStrategyType})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithStrategy, ocsAppsv1.DeploymentStrategy{Type: deploymentConfigStrategyType})
	s.addReplicationControllerWithReplicas(replicationController, 2)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				StrategyTypeRegex: strategyRegex,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentWithStrategy:       {},
				daemonSetWithStrategy:        {},
				deploymentConfigWithStrategy: {},
				replicationController:        {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *UpgradeConfigTestSuite) TestMinPodsUnavailableInteger() {
	const (
		noExplicitRollingUpdate                  = "no-explicit-rolling-update"
		deploymentWithRollingUpdate              = "deployment-rolling-update-min-1"
		deploymentConfigWithRollingUpdate        = "deployment-config-rolling-update-min-1"
		deploymentWithValidRollingUpdate         = "deployment-rolling-update-min-2"
		deploymentConfigWithValidRollingUpdate   = "deployment-config-rolling-update-min-2"
		deploymentWithPercentRollingUpdate       = "deployment-rolling-update-min-10%"
		deploymentConfigWithPercentRollingUpdate = "deployment-config-rolling-update-min-10%"
		replicationController                    = "replication-controller-ignore-integrer"
	)
	maxPodsUnavailable := intstr.FromInt(1)
	maxPodsUnavailableValid := intstr.FromInt(2)
	maxPodsUnavailablePercent := intstr.FromString("10%")
	s.addDeploymentWithStrategy(noExplicitRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType})
	s.addDeploymentWithStrategy(deploymentWithRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxUnavailable: &maxPodsUnavailable}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxUnavailable: &maxPodsUnavailable}})
	s.addDeploymentWithStrategy(deploymentWithPercentRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxUnavailable: &maxPodsUnavailablePercent}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithPercentRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxUnavailable: &maxPodsUnavailablePercent}})
	s.addDeploymentWithStrategy(deploymentWithValidRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxUnavailable: &maxPodsUnavailableValid}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithValidRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxUnavailable: &maxPodsUnavailableValid}})
	s.addReplicationControllerWithReplicas(replicationController, 2)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				StrategyTypeRegex:  "^(RollingUpdate|Rolling)$",
				MinPodsUnavailable: "2",
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				noExplicitRollingUpdate: {
					{Message: "object has no rolling update parameters defined"},
				},
				deploymentWithRollingUpdate: {
					{Message: "object has a max unavailable of 1 but at least 2 is required"},
				},
				deploymentConfigWithRollingUpdate: {
					{Message: "object has a max unavailable of 1 but at least 2 is required"},
				},
				deploymentWithPercentRollingUpdate: {
					{Message: "object has a max unavailable of 10% but at least 2 is required"},
				},
				deploymentConfigWithPercentRollingUpdate: {
					{Message: "object has a max unavailable of 10% but at least 2 is required"},
				},
				deploymentWithValidRollingUpdate:       {},
				deploymentConfigWithValidRollingUpdate: {},
				replicationController:                  {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *UpgradeConfigTestSuite) TestMinPodsUnavailablePercent() {
	const (
		deploymentWithRollingUpdate              = "deployment-rolling-update-min-1"
		deploymentConfigWithRollingUpdate        = "deployment-config-rolling-update-min-1"
		deploymentWithValidRollingUpdate         = "deployment-rolling-update-min-10%"
		deploymentConfigWithValidRollingUpdate   = "deployment-config-rolling-update-min-10%"
		deploymentWithPercentRollingUpdate       = "deployment-rolling-update-min-5%"
		deploymentConfigWithPercentRollingUpdate = "deployment-config-rolling-update-min-5%"
		replicationController                    = "replication-controller-ignore-percent"
	)
	maxPodsUnavailable := intstr.FromInt(1)
	maxPodsUnavailableValid := intstr.FromString("10%")
	maxPodsUnavailablePercent := intstr.FromString("5%")
	s.addDeploymentWithStrategy(deploymentWithRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxUnavailable: &maxPodsUnavailable}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxUnavailable: &maxPodsUnavailable}})
	s.addDeploymentWithStrategy(deploymentWithPercentRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxUnavailable: &maxPodsUnavailablePercent}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithPercentRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxUnavailable: &maxPodsUnavailablePercent}})
	s.addDeploymentWithStrategy(deploymentWithValidRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxUnavailable: &maxPodsUnavailableValid}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithValidRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxUnavailable: &maxPodsUnavailableValid}})
	s.addReplicationControllerWithReplicas(replicationController, 2)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				StrategyTypeRegex:  "^(RollingUpdate|Rolling)$",
				MinPodsUnavailable: "10%",
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentWithRollingUpdate: {
					{Message: "object has a max unavailable of 1 but at least 10% is required"},
				},
				deploymentConfigWithRollingUpdate: {
					{Message: "object has a max unavailable of 1 but at least 10% is required"},
				},
				deploymentWithPercentRollingUpdate: {
					{Message: "object has a max unavailable of 5% but at least 10% is required"},
				},
				deploymentConfigWithPercentRollingUpdate: {
					{Message: "object has a max unavailable of 5% but at least 10% is required"},
				},
				deploymentWithValidRollingUpdate:       {},
				deploymentConfigWithValidRollingUpdate: {},
				replicationController:                  {},
			},
			ExpectInstantiationError: false,
		},
	})
}
func (s *UpgradeConfigTestSuite) TestMaxPodsUnavailableInteger() {
	const (
		noExplicitRollingUpdate                  = "no-explicit-rolling-update"
		deploymentWithRollingUpdate              = "deployment-rolling-update-max-15"
		deploymentConfigWithRollingUpdate        = "deployment-config-rolling-update-max-15"
		deploymentWithValidRollingUpdate         = "deployment-rolling-update-max-10"
		deploymentConfigWithValidRollingUpdate   = "deployment-config-rolling-update-max-10"
		deploymentWithPercentRollingUpdate       = "deployment-rolling-update-max-50%"
		deploymentConfigWithPercentRollingUpdate = "deployment-config-rolling-update-max-50%"
		replicationController                    = "replication-controller-ignore-integer"
	)
	maxPodsUnavailable := intstr.FromInt(15)
	maxPodsUnavailableValid := intstr.FromInt(10)
	maxPodsUnavailablePercent := intstr.FromString("50%")
	s.addDeploymentWithStrategy(noExplicitRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType})
	s.addDeploymentWithStrategy(deploymentWithRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxUnavailable: &maxPodsUnavailable}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxUnavailable: &maxPodsUnavailable}})
	s.addDeploymentWithStrategy(deploymentWithPercentRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxUnavailable: &maxPodsUnavailablePercent}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithPercentRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxUnavailable: &maxPodsUnavailablePercent}})
	s.addDeploymentWithStrategy(deploymentWithValidRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxUnavailable: &maxPodsUnavailableValid}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithValidRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxUnavailable: &maxPodsUnavailableValid}})
	s.addReplicationControllerWithReplicas(replicationController, 2)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				StrategyTypeRegex:  "^(RollingUpdate|Rolling)$",
				MaxPodsUnavailable: "10",
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				noExplicitRollingUpdate: {
					{Message: "object has no rolling update parameters defined"},
				},
				deploymentWithRollingUpdate: {
					{Message: "object has a max unavailable of 15 but no more than 10 is required"},
				},
				deploymentConfigWithRollingUpdate: {
					{Message: "object has a max unavailable of 15 but no more than 10 is required"},
				},
				deploymentWithPercentRollingUpdate: {
					{Message: "object has a max unavailable of 50% but no more than 10 is required"},
				},
				deploymentConfigWithPercentRollingUpdate: {
					{Message: "object has a max unavailable of 50% but no more than 10 is required"},
				},
				deploymentWithValidRollingUpdate:       {},
				deploymentConfigWithValidRollingUpdate: {},
				replicationController:                  {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *UpgradeConfigTestSuite) TestMaxPodsUnavailablePercent() {
	const (
		deploymentWithRollingUpdate              = "deployment-rolling-update-max-5"
		deploymentConfigWithRollingUpdate        = "deployment-config-rolling-update-max-5"
		deploymentWithValidRollingUpdate         = "deployment-rolling-update-max-10%"
		deploymentConfigWithValidRollingUpdate   = "deployment-config-rolling-update-max-10%"
		deploymentWithPercentRollingUpdate       = "deployment-rolling-update-max-50%"
		deploymentConfigWithPercentRollingUpdate = "deployment-config-rolling-update-max-50%"
		replicationController                    = "replication-controller-ignore-percent"
	)
	maxPodsUnavailable := intstr.FromInt(5)
	maxPodsUnavailableValid := intstr.FromString("10%")
	maxPodsUnavailablePercent := intstr.FromString("50%")
	s.addDeploymentWithStrategy(deploymentWithRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxUnavailable: &maxPodsUnavailable}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxUnavailable: &maxPodsUnavailable}})
	s.addDeploymentWithStrategy(deploymentWithPercentRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxUnavailable: &maxPodsUnavailablePercent}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithPercentRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxUnavailable: &maxPodsUnavailablePercent}})
	s.addDeploymentWithStrategy(deploymentWithValidRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxUnavailable: &maxPodsUnavailableValid}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithValidRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxUnavailable: &maxPodsUnavailableValid}})
	s.addReplicationControllerWithReplicas(replicationController, 2)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				StrategyTypeRegex:  "^(RollingUpdate|Rolling)$",
				MaxPodsUnavailable: "10%",
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentWithRollingUpdate: {
					{Message: "object has a max unavailable of 5 but no more than 10% is required"},
				},
				deploymentConfigWithRollingUpdate: {
					{Message: "object has a max unavailable of 5 but no more than 10% is required"},
				},
				deploymentWithPercentRollingUpdate: {
					{Message: "object has a max unavailable of 50% but no more than 10% is required"},
				},
				deploymentConfigWithPercentRollingUpdate: {
					{Message: "object has a max unavailable of 50% but no more than 10% is required"},
				},
				deploymentWithValidRollingUpdate:       {},
				deploymentConfigWithValidRollingUpdate: {},
				replicationController:                  {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *UpgradeConfigTestSuite) TestMaxMinPodsUnavailableInteger() {
	const (
		noExplicitRollingUpdate                  = "no-explicit-rolling-update"
		deploymentWithRollingUpdate              = "deployment-rolling-update-max-15"
		deploymentConfigWithRollingUpdate        = "deployment-config-rolling-update-max-15"
		deploymentWithValidRollingUpdate         = "deployment-rolling-update-max-10"
		deploymentConfigWithValidRollingUpdate   = "deployment-config-rolling-update-max-10"
		deploymentWithPercentRollingUpdate       = "deployment-rolling-update-max-50%"
		deploymentConfigWithPercentRollingUpdate = "deployment-config-rolling-update-max-50%"
		replicationController                    = "replication-controller-ignore-integer"
	)
	maxPodsUnavailable := intstr.FromInt(15)
	maxPodsUnavailableValid := intstr.FromInt(10)
	maxPodsUnavailablePercent := intstr.FromString("50%")
	s.addDeploymentWithStrategy(noExplicitRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType})
	s.addDeploymentWithStrategy(deploymentWithRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxUnavailable: &maxPodsUnavailable}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxUnavailable: &maxPodsUnavailable}})
	s.addDeploymentWithStrategy(deploymentWithPercentRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxUnavailable: &maxPodsUnavailablePercent}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithPercentRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxUnavailable: &maxPodsUnavailablePercent}})
	s.addDeploymentWithStrategy(deploymentWithValidRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxUnavailable: &maxPodsUnavailableValid}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithValidRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxUnavailable: &maxPodsUnavailableValid}})
	s.addReplicationControllerWithReplicas(replicationController, 2)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				StrategyTypeRegex:  "^(RollingUpdate|Rolling)$",
				MinPodsUnavailable: "1",
				MaxPodsUnavailable: "10",
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				noExplicitRollingUpdate: {
					{Message: "object has no rolling update parameters defined"},
				},
				deploymentWithRollingUpdate: {
					{Message: "object has a max unavailable of 15 but at least 1 and no more than 10 is required"},
				},
				deploymentConfigWithRollingUpdate: {
					{Message: "object has a max unavailable of 15 but at least 1 and no more than 10 is required"},
				},
				deploymentWithPercentRollingUpdate: {
					{Message: "object has a max unavailable of 50% but at least 1 and no more than 10 is required"},
				},
				deploymentConfigWithPercentRollingUpdate: {
					{Message: "object has a max unavailable of 50% but at least 1 and no more than 10 is required"},
				},
				deploymentWithValidRollingUpdate:       {},
				deploymentConfigWithValidRollingUpdate: {},
				replicationController:                  {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *UpgradeConfigTestSuite) TestMaxMinPodsUnavailablePercent() {
	const (
		deploymentWithRollingUpdate              = "deployment-rolling-update-max-15"
		deploymentConfigWithRollingUpdate        = "deployment-config-rolling-update-max-15"
		deploymentWithValidRollingUpdate         = "deployment-rolling-update-max-15%"
		deploymentConfigWithValidRollingUpdate   = "deployment-config-rolling-update-max-15%"
		deploymentWithPercentRollingUpdate       = "deployment-rolling-update-max-50%"
		deploymentConfigWithPercentRollingUpdate = "deployment-config-rolling-update-max-50%"
		replicationController                    = "replication-controller-ignore-percent"
	)
	maxPodsUnavailable := intstr.FromInt(15)
	maxPodsUnavailableValid := intstr.FromString("15%")
	maxPodsUnavailablePercent := intstr.FromString("50%")
	s.addDeploymentWithStrategy(deploymentWithRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxUnavailable: &maxPodsUnavailable}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxUnavailable: &maxPodsUnavailable}})
	s.addDeploymentWithStrategy(deploymentWithPercentRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxUnavailable: &maxPodsUnavailablePercent}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithPercentRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxUnavailable: &maxPodsUnavailablePercent}})
	s.addDeploymentWithStrategy(deploymentWithValidRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxUnavailable: &maxPodsUnavailableValid}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithValidRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxUnavailable: &maxPodsUnavailableValid}})
	s.addReplicationControllerWithReplicas(replicationController, 2)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				StrategyTypeRegex:  "^(RollingUpdate|Rolling)$",
				MinPodsUnavailable: "10%",
				MaxPodsUnavailable: "40%",
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentWithRollingUpdate: {
					{Message: "object has a max unavailable of 15 but at least 10% and no more than 40% is required"},
				},
				deploymentConfigWithRollingUpdate: {
					{Message: "object has a max unavailable of 15 but at least 10% and no more than 40% is required"},
				},
				deploymentWithPercentRollingUpdate: {
					{Message: "object has a max unavailable of 50% but at least 10% and no more than 40% is required"},
				},
				deploymentConfigWithPercentRollingUpdate: {
					{Message: "object has a max unavailable of 50% but at least 10% and no more than 40% is required"},
				},
				deploymentWithValidRollingUpdate:       {},
				deploymentConfigWithValidRollingUpdate: {},
				replicationController:                  {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *UpgradeConfigTestSuite) TestMinSurgeInteger() {
	const (
		noExplicitRollingUpdate                  = "no-explicit-rolling-update"
		deploymentInvalidStrategy                = "deployment-invalid-strategy"
		deploymentWithRollingUpdate              = "deployment-rolling-update-min-1"
		deploymentConfigWithRollingUpdate        = "deployment-config-rolling-update-min-1"
		deploymentWithValidRollingUpdate         = "deployment-rolling-update-min-2"
		deploymentConfigWithValidRollingUpdate   = "deployment-config-rolling-update-min-2"
		deploymentWithPercentRollingUpdate       = "deployment-rolling-update-min-10%"
		deploymentConfigWithPercentRollingUpdate = "deployment-config-rolling-update-min-10%"
		daemonSetInvalidStrategy                 = "daemon-set-invalid-strategy"
		daemonSetSurgeTooLow                     = "daemon-set-rolling-update-1"
		daemonSetSurgeValid                      = "daemon-set-rolling-update-2"
		replicationController                    = "replication-controller-ignore-integer"
	)
	maxSurge := intstr.FromInt(1)
	maxSurgeValid := intstr.FromInt(2)
	maxSurgePercent := intstr.FromString("10%")
	s.addDeploymentWithStrategy(noExplicitRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType})
	s.addDeploymentWithStrategy(deploymentInvalidStrategy, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxUnavailable: &maxSurge}})
	s.addDeploymentWithStrategy(deploymentWithRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxSurge: &maxSurge}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxSurge: &maxSurge}})
	s.addDeploymentWithStrategy(deploymentWithPercentRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxSurge: &maxSurgePercent}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithPercentRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxSurge: &maxSurgePercent}})
	s.addDeploymentWithStrategy(deploymentWithValidRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxSurge: &maxSurgeValid}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithValidRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxSurge: &maxSurgeValid}})
	s.addDaemonSetWithStrategy(daemonSetInvalidStrategy, appsv1.DaemonSetUpdateStrategy{Type: appsv1.RollingUpdateDaemonSetStrategyType, RollingUpdate: &appsv1.RollingUpdateDaemonSet{MaxUnavailable: &maxSurge}})
	s.addDaemonSetWithStrategy(daemonSetSurgeTooLow, appsv1.DaemonSetUpdateStrategy{Type: appsv1.RollingUpdateDaemonSetStrategyType, RollingUpdate: &appsv1.RollingUpdateDaemonSet{MaxSurge: &maxSurge}})
	s.addDaemonSetWithStrategy(daemonSetSurgeValid, appsv1.DaemonSetUpdateStrategy{Type: appsv1.RollingUpdateDaemonSetStrategyType, RollingUpdate: &appsv1.RollingUpdateDaemonSet{MaxSurge: &maxSurgeValid}})
	s.addReplicationControllerWithReplicas(replicationController, 2)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				StrategyTypeRegex: "^(RollingUpdate|Rolling)$",
				MinSurge:          "2",
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				noExplicitRollingUpdate: {
					{Message: "object has no rolling update parameters defined"},
				},
				deploymentInvalidStrategy: {
					{Message: "object has a max surge of <nil> but at least 2 is required"},
				},
				deploymentWithRollingUpdate: {
					{Message: "object has a max surge of 1 but at least 2 is required"},
				},
				deploymentConfigWithRollingUpdate: {
					{Message: "object has a max surge of 1 but at least 2 is required"},
				},
				deploymentWithPercentRollingUpdate: {
					{Message: "object has a max surge of 10% but at least 2 is required"},
				},
				deploymentConfigWithPercentRollingUpdate: {
					{Message: "object has a max surge of 10% but at least 2 is required"},
				},
				deploymentWithValidRollingUpdate:       {},
				deploymentConfigWithValidRollingUpdate: {},
				daemonSetInvalidStrategy: {
					{Message: "object has a max surge of <nil> but at least 2 is required"},
				},
				daemonSetSurgeTooLow: {
					{Message: "object has a max surge of 1 but at least 2 is required"},
				},
				daemonSetSurgeValid:   {},
				replicationController: {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *UpgradeConfigTestSuite) TestMinSurgePercent() {
	const (
		deploymentWithRollingUpdate              = "deployment-rolling-update-min-1"
		deploymentConfigWithRollingUpdate        = "deployment-config-rolling-update-min-1"
		deploymentWithValidRollingUpdate         = "deployment-rolling-update-min-10%"
		deploymentConfigWithValidRollingUpdate   = "deployment-config-rolling-update-min-10%"
		deploymentWithPercentRollingUpdate       = "deployment-rolling-update-min-5%"
		deploymentConfigWithPercentRollingUpdate = "deployment-config-rolling-update-min-5%"
		replicationController                    = "replication-controller-ignore-percent"
	)
	maxSurge := intstr.FromInt(1)
	maxSurgeValid := intstr.FromString("10%")
	maxSurgePercent := intstr.FromString("5%")
	s.addDeploymentWithStrategy(deploymentWithRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxSurge: &maxSurge}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxSurge: &maxSurge}})
	s.addDeploymentWithStrategy(deploymentWithPercentRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxSurge: &maxSurgePercent}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithPercentRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxSurge: &maxSurgePercent}})
	s.addDeploymentWithStrategy(deploymentWithValidRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxSurge: &maxSurgeValid}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithValidRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxSurge: &maxSurgeValid}})
	s.addReplicationControllerWithReplicas(replicationController, 2)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				StrategyTypeRegex: "^(RollingUpdate|Rolling)$",
				MinSurge:          "10%",
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentWithRollingUpdate: {
					{Message: "object has a max surge of 1 but at least 10% is required"},
				},
				deploymentConfigWithRollingUpdate: {
					{Message: "object has a max surge of 1 but at least 10% is required"},
				},
				deploymentWithPercentRollingUpdate: {
					{Message: "object has a max surge of 5% but at least 10% is required"},
				},
				deploymentConfigWithPercentRollingUpdate: {
					{Message: "object has a max surge of 5% but at least 10% is required"},
				},
				deploymentWithValidRollingUpdate:       {},
				deploymentConfigWithValidRollingUpdate: {},
				replicationController:                  {},
			},
			ExpectInstantiationError: false,
		},
	})
}
func (s *UpgradeConfigTestSuite) TestMaxSurgeInteger() {
	const (
		noExplicitRollingUpdate                  = "no-explicit-rolling-update"
		deploymentWithRollingUpdate              = "deployment-rolling-update-max-15"
		deploymentConfigWithRollingUpdate        = "deployment-config-rolling-update-max-15"
		deploymentWithValidRollingUpdate         = "deployment-rolling-update-max-10"
		deploymentConfigWithValidRollingUpdate   = "deployment-config-rolling-update-max-10"
		deploymentWithPercentRollingUpdate       = "deployment-rolling-update-max-50%"
		deploymentConfigWithPercentRollingUpdate = "deployment-config-rolling-update-max-50%"
		replicationController                    = "replication-controller-ignore-integer"
	)
	maxSurge := intstr.FromInt(15)
	maxSurgeValid := intstr.FromInt(10)
	maxSurgePercent := intstr.FromString("50%")
	s.addDeploymentWithStrategy(noExplicitRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType})
	s.addDeploymentWithStrategy(deploymentWithRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxSurge: &maxSurge}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxSurge: &maxSurge}})
	s.addDeploymentWithStrategy(deploymentWithPercentRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxSurge: &maxSurgePercent}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithPercentRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxSurge: &maxSurgePercent}})
	s.addDeploymentWithStrategy(deploymentWithValidRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxSurge: &maxSurgeValid}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithValidRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxSurge: &maxSurgeValid}})
	s.addReplicationControllerWithReplicas(replicationController, 2)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				StrategyTypeRegex: "^(RollingUpdate|Rolling)$",
				MaxSurge:          "10",
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				noExplicitRollingUpdate: {
					{Message: "object has no rolling update parameters defined"},
				},
				deploymentWithRollingUpdate: {
					{Message: "object has a max surge of 15 but no more than 10 is required"},
				},
				deploymentConfigWithRollingUpdate: {
					{Message: "object has a max surge of 15 but no more than 10 is required"},
				},
				deploymentWithPercentRollingUpdate: {
					{Message: "object has a max surge of 50% but no more than 10 is required"},
				},
				deploymentConfigWithPercentRollingUpdate: {
					{Message: "object has a max surge of 50% but no more than 10 is required"},
				},
				deploymentWithValidRollingUpdate:       {},
				deploymentConfigWithValidRollingUpdate: {},
				replicationController:                  {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *UpgradeConfigTestSuite) TestMaxSurgePercent() {
	const (
		deploymentWithRollingUpdate              = "deployment-rolling-update-max-5"
		deploymentConfigWithRollingUpdate        = "deployment-config-rolling-update-max-5"
		deploymentWithValidRollingUpdate         = "deployment-rolling-update-max-10%"
		deploymentConfigWithValidRollingUpdate   = "deployment-config-rolling-update-max-10%"
		deploymentWithPercentRollingUpdate       = "deployment-rolling-update-max-50%"
		deploymentConfigWithPercentRollingUpdate = "deployment-config-rolling-update-max-50%"
		replicationController                    = "replication-controller-ignore-percent"
	)
	maxSurge := intstr.FromInt(5)
	maxSurgeValid := intstr.FromString("10%")
	maxSurgePercent := intstr.FromString("50%")
	s.addDeploymentWithStrategy(deploymentWithRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxSurge: &maxSurge}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxSurge: &maxSurge}})
	s.addDeploymentWithStrategy(deploymentWithPercentRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxSurge: &maxSurgePercent}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithPercentRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxSurge: &maxSurgePercent}})
	s.addDeploymentWithStrategy(deploymentWithValidRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxSurge: &maxSurgeValid}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithValidRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxSurge: &maxSurgeValid}})
	s.addReplicationControllerWithReplicas(replicationController, 2)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				StrategyTypeRegex: "^(RollingUpdate|Rolling)$",
				MaxSurge:          "10%",
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentWithRollingUpdate: {
					{Message: "object has a max surge of 5 but no more than 10% is required"},
				},
				deploymentConfigWithRollingUpdate: {
					{Message: "object has a max surge of 5 but no more than 10% is required"},
				},
				deploymentWithPercentRollingUpdate: {
					{Message: "object has a max surge of 50% but no more than 10% is required"},
				},
				deploymentConfigWithPercentRollingUpdate: {
					{Message: "object has a max surge of 50% but no more than 10% is required"},
				},
				deploymentWithValidRollingUpdate:       {},
				deploymentConfigWithValidRollingUpdate: {},
				replicationController:                  {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *UpgradeConfigTestSuite) TestMaxMinSurgeInteger() {
	const (
		noExplicitRollingUpdate                  = "no-explicit-rolling-update"
		deploymentWithRollingUpdate              = "deployment-rolling-update-max-15"
		deploymentConfigWithRollingUpdate        = "deployment-config-rolling-update-max-15"
		deploymentWithValidRollingUpdate         = "deployment-rolling-update-max-10"
		deploymentConfigWithValidRollingUpdate   = "deployment-config-rolling-update-max-10"
		deploymentWithPercentRollingUpdate       = "deployment-rolling-update-max-50%"
		deploymentConfigWithPercentRollingUpdate = "deployment-config-rolling-update-max-50%"
		replicationController                    = "replication-controller-ignore-integer"
	)
	maxSurge := intstr.FromInt(15)
	maxSurgeValid := intstr.FromInt(10)
	maxSurgePercent := intstr.FromString("50%")
	s.addDeploymentWithStrategy(noExplicitRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType})
	s.addDeploymentWithStrategy(deploymentWithRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxSurge: &maxSurge}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxSurge: &maxSurge}})
	s.addDeploymentWithStrategy(deploymentWithPercentRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxSurge: &maxSurgePercent}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithPercentRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxSurge: &maxSurgePercent}})
	s.addDeploymentWithStrategy(deploymentWithValidRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxSurge: &maxSurgeValid}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithValidRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxSurge: &maxSurgeValid}})
	s.addReplicationControllerWithReplicas(replicationController, 3)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				StrategyTypeRegex: "^(RollingUpdate|Rolling)$",
				MinSurge:          "1",
				MaxSurge:          "10",
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				noExplicitRollingUpdate: {
					{Message: "object has no rolling update parameters defined"},
				},
				deploymentWithRollingUpdate: {
					{Message: "object has a max surge of 15 but at least 1 and no more than 10 is required"},
				},
				deploymentConfigWithRollingUpdate: {
					{Message: "object has a max surge of 15 but at least 1 and no more than 10 is required"},
				},
				deploymentWithPercentRollingUpdate: {
					{Message: "object has a max surge of 50% but at least 1 and no more than 10 is required"},
				},
				deploymentConfigWithPercentRollingUpdate: {
					{Message: "object has a max surge of 50% but at least 1 and no more than 10 is required"},
				},
				deploymentWithValidRollingUpdate:       {},
				deploymentConfigWithValidRollingUpdate: {},
				replicationController:                  {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *UpgradeConfigTestSuite) TestMaxMinSurgePercent() {
	const (
		deploymentWithRollingUpdate              = "deployment-rolling-update-max-15"
		deploymentConfigWithRollingUpdate        = "deployment-config-rolling-update-max-15"
		deploymentWithValidRollingUpdate         = "deployment-rolling-update-max-15%"
		deploymentConfigWithValidRollingUpdate   = "deployment-config-rolling-update-max-15%"
		deploymentWithPercentRollingUpdate       = "deployment-rolling-update-max-50%"
		deploymentConfigWithPercentRollingUpdate = "deployment-config-rolling-update-max-50%"
		replicationController                    = "replication-controller-ignore-percent"
	)
	maxSurge := intstr.FromInt(15)
	maxSurgeValid := intstr.FromString("15%")
	maxSurgePercent := intstr.FromString("50%")
	s.addDeploymentWithStrategy(deploymentWithRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxSurge: &maxSurge}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxSurge: &maxSurge}})
	s.addDeploymentWithStrategy(deploymentWithPercentRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxSurge: &maxSurgePercent}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithPercentRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxSurge: &maxSurgePercent}})
	s.addDeploymentWithStrategy(deploymentWithValidRollingUpdate, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxSurge: &maxSurgeValid}})
	s.addDeploymentConfigWithStrategy(deploymentConfigWithValidRollingUpdate, ocsAppsv1.DeploymentStrategy{Type: ocsAppsv1.DeploymentStrategyTypeRolling, RollingParams: &ocsAppsv1.RollingDeploymentStrategyParams{MaxSurge: &maxSurgeValid}})
	s.addReplicationControllerWithReplicas(replicationController, 2)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				StrategyTypeRegex: "^(RollingUpdate|Rolling)$",
				MinSurge:          "10%",
				MaxSurge:          "40%",
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentWithRollingUpdate: {
					{Message: "object has a max surge of 15 but at least 10% and no more than 40% is required"},
				},
				deploymentConfigWithRollingUpdate: {
					{Message: "object has a max surge of 15 but at least 10% and no more than 40% is required"},
				},
				deploymentWithPercentRollingUpdate: {
					{Message: "object has a max surge of 50% but at least 10% and no more than 40% is required"},
				},
				deploymentConfigWithPercentRollingUpdate: {
					{Message: "object has a max surge of 50% but at least 10% and no more than 40% is required"},
				},
				deploymentWithValidRollingUpdate:       {},
				deploymentConfigWithValidRollingUpdate: {},
				replicationController:                  {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *UpgradeConfigTestSuite) TestTemplateConfig() {
	const (
		validDeployment = "valid-deployment"
	)
	maxSurgeValid := intstr.FromString("10%")
	maxUnavailableValid := intstr.FromInt(10)
	s.addDeploymentWithStrategy(validDeployment, appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType, RollingUpdate: &appsv1.RollingUpdateDeployment{MaxUnavailable: &maxUnavailableValid, MaxSurge: &maxSurgeValid}})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				StrategyTypeRegex:  "^(RollingUpdate|Rolling)$",
				MinPodsUnavailable: "10",
				MaxPodsUnavailable: "40",
				MinSurge:           "10%",
				MaxSurge:           "40%",
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				validDeployment: {},
			},
			ExpectInstantiationError: false,
		},
		{
			Param: params.Params{
				StrategyTypeRegex:  "^(RollingUpdate|Rolling)$",
				MinPodsUnavailable: "-1",
			},
			ExpectInstantiationError: true,
		},
		{
			Param: params.Params{
				StrategyTypeRegex:  "^(RollingUpdate|Rolling)$",
				MaxPodsUnavailable: "1%1",
			},
			ExpectInstantiationError: true,
		},
		{
			Param: params.Params{
				StrategyTypeRegex: "^(RollingUpdate|Rolling)$",
				MinSurge:          "115%",
			},
			ExpectInstantiationError: true,
		},
		{
			Param: params.Params{
				StrategyTypeRegex: "^(RollingUpdate|Rolling)$",
				MaxSurge:          "1%1%",
			},
			ExpectInstantiationError: true,
		},
	})
}
