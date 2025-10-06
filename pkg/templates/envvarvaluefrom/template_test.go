package envvarvaluefrom

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"golang.stackrox.io/kube-linter/internal/pointers"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/envvarvaluefrom/internal/params"
	coreV1 "k8s.io/api/core/v1"
)

const (
	targetDeploymentName = "deployment-1"
)

func TestEnvVarValueFrom(t *testing.T) {
	suite.Run(t, new(EnVarValueFromTestSuite))
}

type EnVarValueFromTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *EnVarValueFromTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

type sourceReference struct {
	Name     string
	Key      string
	Optional *bool
}

type envReference struct {
	Name   string
	Kind   string
	Source sourceReference
}

func makeSecretSource(descriptor sourceReference) *coreV1.EnvVarSource {
	return &coreV1.EnvVarSource{
		SecretKeyRef: &coreV1.SecretKeySelector{
			LocalObjectReference: coreV1.LocalObjectReference{
				Name: descriptor.Name,
			},
			Key:      descriptor.Key,
			Optional: descriptor.Optional,
		},
	}
}

func makeConfigMapSource(descriptor sourceReference) *coreV1.EnvVarSource {
	return &coreV1.EnvVarSource{
		ConfigMapKeyRef: &coreV1.ConfigMapKeySelector{
			LocalObjectReference: coreV1.LocalObjectReference{
				Name: descriptor.Name,
			},
			Key:      descriptor.Key,
			Optional: descriptor.Optional,
		},
	}
}

func (s *EnVarValueFromTestSuite) addContainerWithEnvFromSecret(name string, envRef envReference) {
	var valueFrom *coreV1.EnvVarSource
	switch envRef.Kind {
	case "secret":
		valueFrom = makeSecretSource(envRef.Source)
	case "configmap":
		valueFrom = makeConfigMapSource(envRef.Source)
	default:
		require.FailNow(s.T(), fmt.Sprintf("Unknown source kind %s", envRef.Kind))
	}

	s.ctx.AddContainerToDeployment(s.T(), name, coreV1.Container{
		Name: "container",
		Env: []coreV1.EnvVar{
			{
				Name:      "ENV_1",
				ValueFrom: valueFrom,
			},
		},
	})
}

func (s *EnVarValueFromTestSuite) TestDeploymentWithoutEnvPasses() {
	s.ctx.AddMockDeployment(s.T(), targetDeploymentName)

	s.ctx.AddContainerToDeployment(s.T(), targetDeploymentName, coreV1.Container{
		Name: "container",
		Env:  []coreV1.EnvVar{},
	})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				targetDeploymentName: {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *EnVarValueFromTestSuite) TestDeploymentWithDirectEnvPasses() {
	s.ctx.AddMockDeployment(s.T(), targetDeploymentName)

	s.ctx.AddContainerToDeployment(s.T(), targetDeploymentName, coreV1.Container{
		Name: "container",
		Env: []coreV1.EnvVar{
			{
				Name:  "ENV_1",
				Value: "Value",
			},
		},
	})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				targetDeploymentName: {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *EnVarValueFromTestSuite) TestDeploymentWithUnknownSecret() {
	s.ctx.AddMockDeployment(s.T(), targetDeploymentName)

	s.addContainerWithEnvFromSecret(targetDeploymentName, envReference{
		Name: "my-secret",
		Kind: "secret",
		Source: sourceReference{
			Name:     "foo",
			Key:      "bar",
			Optional: pointers.Bool(false),
		},
	})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				targetDeploymentName: {{
					Message: "The container \"container\" is referring to an unknown secret \"foo\"",
				}},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *EnVarValueFromTestSuite) TestDeploymentWithNoOptionalSecret() {
	s.ctx.AddMockDeployment(s.T(), targetDeploymentName)

	s.addContainerWithEnvFromSecret(targetDeploymentName, envReference{
		Name: "my-secret",
		Kind: "secret",
		Source: sourceReference{
			Name:     "foo",
			Key:      "bar",
			Optional: nil,
		},
	})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				targetDeploymentName: {{
					Message: "The container \"container\" is referring to an unknown secret \"foo\"",
				}},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *EnVarValueFromTestSuite) TestDeploymentWithUnknownOptionalSecretPasses() {
	s.ctx.AddMockDeployment(s.T(), targetDeploymentName)

	s.addContainerWithEnvFromSecret(targetDeploymentName, envReference{
		Name: "my-secret",
		Kind: "secret",
		Source: sourceReference{
			Name:     "foo",
			Key:      "bar",
			Optional: pointers.Bool(true),
		},
	})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				targetDeploymentName: {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *EnVarValueFromTestSuite) TestDeploymentWithUnknownConfigMap() {
	s.ctx.AddMockDeployment(s.T(), targetDeploymentName)

	s.addContainerWithEnvFromSecret(targetDeploymentName, envReference{
		Name: "my_config_var",
		Kind: "configmap",
		Source: sourceReference{
			Name:     "foo",
			Key:      "bar",
			Optional: pointers.Bool(false),
		},
	})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				targetDeploymentName: {{
					Message: "The container \"container\" is referring to an unknown config map \"foo\"",
				}},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *EnVarValueFromTestSuite) TestDeploymentWithUnknownOptionalConfigMapPasses() {
	s.ctx.AddMockDeployment(s.T(), targetDeploymentName)

	s.addContainerWithEnvFromSecret(targetDeploymentName, envReference{
		Name: "my_config_var",
		Kind: "configmap",
		Source: sourceReference{
			Name:     "foo",
			Key:      "bar",
			Optional: pointers.Bool(true),
		},
	})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				targetDeploymentName: {},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *EnVarValueFromTestSuite) TestDeploymentWithNoOptionalConfigMap() {
	s.ctx.AddMockDeployment(s.T(), targetDeploymentName)

	s.addContainerWithEnvFromSecret(targetDeploymentName, envReference{
		Name: "my-config",
		Kind: "configmap",
		Source: sourceReference{
			Name:     "foo",
			Key:      "bar",
			Optional: nil,
		},
	})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				targetDeploymentName: {{
					Message: "The container \"container\" is referring to an unknown config map \"foo\"",
				}},
			},
			ExpectInstantiationError: false,
		},
	})
}
