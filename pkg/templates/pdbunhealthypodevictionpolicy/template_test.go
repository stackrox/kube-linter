package pdbunhealthypodevictionpolicy

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/pdbunhealthypodevictionpolicy/internal/params"
	pdbv1 "k8s.io/api/policy/v1"
)

func TestUnhealthyPodEvictionPolicyPDB(t *testing.T) {
	suite.Run(t, new(UnhealthyPodEvictionPolicyPDBTestSuite))
}

type UnhealthyPodEvictionPolicyPDBTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *UnhealthyPodEvictionPolicyPDBTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *UnhealthyPodEvictionPolicyPDBTestSuite) TestNoUnhealthyPodEvictionPolicy() {
	s.ctx.AddMockPodDisruptionBudget(s.T(), "test-pdb-no-unhealthy-pod-eviction-policy")
	s.ctx.ModifyPodDisruptionBudget(s.T(), "test-pdb-no-unhealthy-pod-eviction-policy", func(pdb *pdbv1.PodDisruptionBudget) {
		pdb.Spec.UnhealthyPodEvictionPolicy = nil
	})
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				"test-pdb-no-unhealthy-pod-eviction-policy": {{Message: "unhealthyPodEvictionPolicy is not explicitly set"}},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *UnhealthyPodEvictionPolicyPDBTestSuite) TestUnhealthyPodEvictionPolicyIsSet() {
	s.ctx.AddMockPodDisruptionBudget(s.T(), "test-pdb-unhealthy-pod-eviction-policy-is-set")
	s.ctx.ModifyPodDisruptionBudget(s.T(), "test-pdb-unhealthy-pod-eviction-policy-is-set", func(pdb *pdbv1.PodDisruptionBudget) {
		var policy pdbv1.UnhealthyPodEvictionPolicyType = "AlwaysAllow"
		pdb.Spec.UnhealthyPodEvictionPolicy = &policy
	})
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				"test-pdb-unhealthy-pod-eviction-policy-is-set": nil,
			},
			ExpectInstantiationError: false,
		},
	})
}
