package accesstoresources

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/accesstoresources/internal/params"
	rbacV1 "k8s.io/api/rbac/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	role1               = "role1"
	role2               = "role2"
	clusterRole1        = "cluster-role1"
	clusterRole2        = "cluster-role2"
	roleBinding1        = "role-binding1"
	roleBinding2        = "role-binding2"
	clusterRoleBinding1 = "cluster-role-binding1"
	clusterRoleBinding2 = "cluster-role-binding2"
	namespace1          = "namespace-dev"
	namespace2          = "namespace-test"
)

var rules1 = []rbacV1.PolicyRule{{
	APIGroups: []string{""},
	Resources: []string{"deployments"},
	Verbs:     []string{"get", "list", "watch", "create", "delete"},
}}
var rules2 = []rbacV1.PolicyRule{{
	APIGroups: []string{""},
	Resources: []string{"services"},
	Verbs:     []string{"*"},
}}
var rules3 = []rbacV1.PolicyRule{{
	APIGroups: []string{""},
	Resources: []string{"secrets", "configmap"},
	Verbs:     []string{"get", "list", "watch", "create", "delete"},
}}
var rules4 = []rbacV1.PolicyRule{{
	APIGroups: []string{""},
	Resources: []string{"*"},
	Verbs:     []string{"get", "list", "watch"},
}}
var labels1 = map[string]string{"app": "test1", "location": "south"}
var aggRule = rbacV1.AggregationRule{
	ClusterRoleSelectors: []metaV1.LabelSelector{{
		MatchLabels: labels1,
	}},
}
var roleRef1 = rbacV1.RoleRef{
	APIGroup: rbacV1.GroupName,
	Kind:     objectkinds.Role,
	Name:     role1,
}
var roleRef2 = rbacV1.RoleRef{
	APIGroup: rbacV1.GroupName,
	Kind:     objectkinds.Role,
	Name:     role2,
}
var clusterRoleRef1 = rbacV1.RoleRef{
	APIGroup: rbacV1.GroupName,
	Kind:     objectkinds.ClusterRole,
	Name:     clusterRole1,
}
var clusterRoleRef2 = rbacV1.RoleRef{
	APIGroup: rbacV1.GroupName,
	Kind:     objectkinds.ClusterRole,
	Name:     clusterRole2,
}
var subjects1 = []rbacV1.Subject{{
	Kind:      "ServiceAccount",
	APIGroup:  "",
	Name:      "account1",
	Namespace: namespace1,
}}
var subjects2 = []rbacV1.Subject{{
	Kind:      "ServiceAccount",
	APIGroup:  "",
	Name:      "account2",
	Namespace: namespace2,
}}

func TestAccessToSecrets(t *testing.T) {
	suite.Run(t, new(AccessToSecretsTestSuite))
}

type AccessToSecretsTestSuite struct {
	templates.TemplateTestSuite
	ctx *mocks.MockLintContext
}

func (s *AccessToSecretsTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *AccessToSecretsTestSuite) addRole(name, namespace string, rules []rbacV1.PolicyRule) {
	s.ctx.AddMockRole(s.T(), name, namespace)
	s.ctx.ModifyRole(s.T(), name, func(role *rbacV1.Role) {
		role.Rules = append(role.Rules, rules...)
	})
}

func (s *AccessToSecretsTestSuite) addClusterRole(name string, labels map[string]string, rules []rbacV1.PolicyRule, aggRule rbacV1.AggregationRule) {
	s.ctx.AddMockClusterRole(s.T(), name)
	s.ctx.ModifyClusterRole(s.T(), name, func(crole *rbacV1.ClusterRole) {
		crole.Labels = labels
		crole.Rules = append(crole.Rules, rules...)
		crole.AggregationRule = &aggRule
	})
}

func (s *AccessToSecretsTestSuite) addRoleBinding(name, namespace string, ref rbacV1.RoleRef, subjects []rbacV1.Subject) {
	s.ctx.AddMockRoleBinding(s.T(), name, namespace)
	s.ctx.ModifyRoleBinding(s.T(), name, func(rolebinding *rbacV1.RoleBinding) {
		rolebinding.RoleRef = ref
		rolebinding.Subjects = append(rolebinding.Subjects, subjects...)
	})
}

func (s *AccessToSecretsTestSuite) addClusterRoleBinding(name string, ref rbacV1.RoleRef, subjects []rbacV1.Subject) {
	s.ctx.AddMockClusterRoleBinding(s.T(), name)
	s.ctx.ModifyClusterRoleBinding(s.T(), name, func(clusterrolebinding *rbacV1.ClusterRoleBinding) {
		clusterrolebinding.RoleRef = ref
		clusterrolebinding.Subjects = append(clusterrolebinding.Subjects, subjects...)
	})
}

func (s *AccessToSecretsTestSuite) TestNoRestrictions() {

	s.addRole(role1, namespace1, rules1)
	s.addRoleBinding(roleBinding1, namespace1, roleRef1, subjects1)
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param:                    params.Params{FlagRolesNotFound: false},
			Diagnostics:              nil,
			ExpectInstantiationError: false,
		},
	})
	s.addRole(role2, namespace2, rules2)
	s.addRoleBinding(roleBinding2, namespace2, roleRef2, subjects2)
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param:                    params.Params{FlagRolesNotFound: false},
			Diagnostics:              nil,
			ExpectInstantiationError: false,
		},
	})
	s.addClusterRole(clusterRole2, nil, rules1, aggRule)
	s.addClusterRoleBinding(clusterRoleBinding2, clusterRoleRef2, subjects1)
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param:                    params.Params{FlagRolesNotFound: false},
			Diagnostics:              nil,
			ExpectInstantiationError: false,
		},
	})
}

func (s *AccessToSecretsTestSuite) TestRoleNotFound() {
	s.addRole(role1, namespace1, rules1)
	s.addRoleBinding(roleBinding1, namespace2, roleRef1, subjects1)
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{FlagRolesNotFound: true},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				roleBinding1: {
					{Message: fmt.Sprintf("role %q in namespace %q not found", role1, namespace2)},
				},
			},
			ExpectInstantiationError: false,
		},
	})

	s.addClusterRole(clusterRole1, nil, rules2, aggRule)
	s.addClusterRoleBinding(clusterRoleBinding1, clusterRoleRef2, subjects1)
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				FlagRolesNotFound: true,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				roleBinding1: {
					{Message: fmt.Sprintf("role %q in namespace %q not found", role1, namespace2)},
				},
				clusterRoleBinding1: {
					{Message: fmt.Sprintf("clusterrole %q not found", clusterRoleRef2.Name)},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *AccessToSecretsTestSuite) TestRoleWithNoAccessToSecrets() {
	rulesForRole := rules1
	rulesForRole = append(rulesForRole, rules3...)
	s.addRole(role1, namespace1, rulesForRole)
	s.addRoleBinding(roleBinding1, namespace1, roleRef1, subjects1)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				FlagRolesNotFound: false,
				Resources:         []string{"^secrets$"},
				Verbs:             []string{"^get$", "^delete$", "^create$"},
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				roleBinding1: {
					{Message: fmt.Sprintf("binding to %q role that has [get create delete] access to [secrets]", role1)},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *AccessToSecretsTestSuite) TestClusterRoleWithNoAccessToSecrets() {
	rulesForClusterRole := rules2
	rulesForClusterRole = append(rulesForClusterRole, rules4...)
	s.addClusterRole(clusterRole1, nil, rulesForClusterRole, aggRule)
	s.addClusterRoleBinding(clusterRoleBinding1, clusterRoleRef1, subjects1)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				FlagRolesNotFound: false,
				Resources:         []string{"^secrets$"},
				Verbs:             []string{"^get$", "^create$", "^watch$"},
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				clusterRoleBinding1: {
					{Message: fmt.Sprintf("binding to %q clusterrole that has [get watch] access to [*]", clusterRole1)},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *AccessToSecretsTestSuite) TestClusterRoleWithAggregationRule() {
	s.addClusterRole(clusterRole1, nil, rules2, aggRule)
	s.addClusterRole(clusterRole2, labels1, rules4, rbacV1.AggregationRule{})
	s.addClusterRoleBinding(clusterRoleBinding1, clusterRoleRef1, subjects1)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				FlagRolesNotFound: false,
				Resources:         []string{"^secrets$"},
				Verbs:             []string{"^get$", "^create$", "^watch$"},
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				clusterRoleBinding1: {
					{Message: fmt.Sprintf("binding via aggregationRule to %q clusterrole that has [get watch] access to [*]", clusterRole2)},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *AccessToSecretsTestSuite) TestRoleWithNoAccessToDeployments() {
	rulesForRole := rules1
	rulesForRole = append(rulesForRole, rules2...)
	s.addRole(role1, namespace1, rulesForRole)
	s.addRoleBinding(roleBinding1, namespace1, roleRef1, subjects1)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				FlagRolesNotFound: false,
				Resources:         []string{"^deployments$", "^cronjobs$"},
				Verbs:             []string{"^get$", "^create$", "^watch$"},
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				roleBinding1: {
					{Message: fmt.Sprintf("binding to %q role that has [get watch create] access to [deployments]", role1)},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}
