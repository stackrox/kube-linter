package dnsconfigoptions

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/dnsconfigoptions/internal/params"
	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

func TestDnsconfigOptions(t *testing.T) {
	suite.Run(t, new(DnsconfigOptionsTestSuite))
}

type DnsconfigOptionsTestSuite struct {
	templates.TemplateTestSuite
	ctx *mocks.MockLintContext
}

const (
	templateKey    = "dnsconfig-options"
	deploymentName = "deployment"
	key            = "ndots"
	value          = "2"
)

func (s *DnsconfigOptionsTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *DnsconfigOptionsTestSuite) TestIgnoreDnsconfigOptionsCheckOnObjectWithoutDnsconfig() {
	s.ctx.AddMockClusterRole(s.T(), deploymentName)
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param:                    params.Params{Key: key, Value: value},
			Diagnostics:              nil,
			ExpectInstantiationError: false,
		},
	})
}

func (s *DnsconfigOptionsTestSuite) TestNoPodTemplateSpecDnsconfigDefined() {
	s.ctx.AddMockDeployment(s.T(), deploymentName)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{Key: key, Value: value},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentName: {
					{
						Message: "Object does not define any DNSConfig rules.",
					},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DnsconfigOptionsTestSuite) TestNoDnsconfigOptionsDefined() {
	s.ctx.AddMockDeployment(s.T(), deploymentName)

	s.ctx.ModifyDeployment(s.T(), deploymentName, func(deployment *appsV1.Deployment) {
		dnsconfig := &v1.PodDNSConfig{
			Options: nil,
		}
		deployment.Spec.Template.Spec.DNSConfig = dnsconfig
	})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{Key: key, Value: value},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentName: {
					{
						Message: "Object does not define any DNSConfig Options.",
					},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DnsconfigOptionsTestSuite) TestDnsconfigOptionsNotMatched() {
	s.ctx.AddMockDeployment(s.T(), deploymentName)

	v := "5"
	s.ctx.ModifyDeployment(s.T(), deploymentName, func(deployment *appsV1.Deployment) {
		dnsconfig := &v1.PodDNSConfig{
			Options: []v1.PodDNSConfigOption{
				{
					Name:  "foo",
					Value: &v,
				},
			},
		}
		deployment.Spec.Template.Spec.DNSConfig = dnsconfig
	})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{Key: key, Value: value},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentName: {
					{
						Message: fmt.Sprintf("DNSConfig Options \"%s:%s\" not found.", key, value),
					},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *DnsconfigOptionsTestSuite) TestDnsconfigOptionsMatched() {
	s.ctx.AddMockDeployment(s.T(), deploymentName)

	v := "2"
	s.ctx.ModifyDeployment(s.T(), deploymentName, func(deployment *appsV1.Deployment) {
		dnsconfig := &v1.PodDNSConfig{
			Options: []v1.PodDNSConfigOption{
				{
					Name:  "ndots",
					Value: &v,
				},
			},
		}
		deployment.Spec.Template.Spec.DNSConfig = dnsconfig
	})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param:                    params.Params{Key: key, Value: value},
			Diagnostics:              nil,
			ExpectInstantiationError: false,
		},
	})
}
