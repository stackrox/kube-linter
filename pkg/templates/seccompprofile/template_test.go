package seccompprofile

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/seccompprofile/internal/params"

	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

const (
	containerName           = "test-container"
	initContainerName       = "test-init-container"
	ephemeralContainerName  = "test-ephemeral-container"
	runtimeDefaultProfile   = "RuntimeDefault"
	localhostProfile        = "Localhost"
	unconfinedProfile       = "Unconfined"
	expectedDiagMessageTmpl = "container %q does not have a seccomp profile set. " +
		"Set seccompProfile.type to \"RuntimeDefault\" or \"Localhost\" in the securityContext of the container or the pod."
)

func TestSeccompProfile(t *testing.T) {
	suite.Run(t, new(SeccompProfileTestSuite))
}

type SeccompProfileTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *SeccompProfileTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func expectedDiag(containerName string) []diagnostic.Diagnostic {
	return []diagnostic.Diagnostic{
		{Message: fmt.Sprintf(expectedDiagMessageTmpl, containerName)},
	}
}

func (s *SeccompProfileTestSuite) addDeployment(name string, podSC *v1.PodSecurityContext, containers, initContainers []v1.Container, ephemeralContainers []v1.EphemeralContainer) {
	s.ctx.AddMockDeployment(s.T(), name)
	if podSC != nil {
		s.ctx.AddSecurityContextToDeployment(s.T(), name, podSC)
	}
	for _, container := range containers {
		s.ctx.AddContainerToDeployment(s.T(), name, container)
	}
	if len(initContainers) > 0 || len(ephemeralContainers) > 0 {
		s.ctx.ModifyDeployment(s.T(), name, func(deployment *appsV1.Deployment) {
			deployment.Spec.Template.Spec.InitContainers = append(deployment.Spec.Template.Spec.InitContainers, initContainers...)
			deployment.Spec.Template.Spec.EphemeralContainers = append(deployment.Spec.Template.Spec.EphemeralContainers, ephemeralContainers...)
		})
	}
}

var localhostProfileRef = "my-profiles/profile.json"

func containerWithProfile(profileType string) v1.Container {
	return v1.Container{
		Name: containerName,
		SecurityContext: &v1.SecurityContext{
			SeccompProfile: &v1.SeccompProfile{Type: v1.SeccompProfileType(profileType)},
		},
	}
}

func (s *SeccompProfileTestSuite) TestNoSeccompProfileFlagged() {
	const (
		noSecurityContext   = "no-security-context"
		emptySecurityCtx    = "empty-security-context"
		unconfinedContainer = "unconfined-container"
		unconfinedPod       = "unconfined-pod"
		emptyProfileType    = "empty-profile-type"
	)

	s.addDeployment(noSecurityContext, nil,
		[]v1.Container{{Name: containerName}}, nil, nil)
	s.addDeployment(emptySecurityCtx, nil,
		[]v1.Container{{Name: containerName, SecurityContext: &v1.SecurityContext{}}}, nil, nil)
	s.addDeployment(unconfinedContainer, nil,
		[]v1.Container{containerWithProfile(unconfinedProfile)}, nil, nil)
	s.addDeployment(unconfinedPod,
		&v1.PodSecurityContext{SeccompProfile: &v1.SeccompProfile{Type: v1.SeccompProfileTypeUnconfined}},
		[]v1.Container{{Name: containerName}}, nil, nil)
	s.addDeployment(emptyProfileType, nil,
		[]v1.Container{containerWithProfile("")}, nil, nil)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				noSecurityContext:   expectedDiag(containerName),
				emptySecurityCtx:    expectedDiag(containerName),
				unconfinedContainer: expectedDiag(containerName),
				unconfinedPod:       expectedDiag(containerName),
				emptyProfileType:    expectedDiag(containerName),
			},
		},
	})
}

func (s *SeccompProfileTestSuite) TestSeccompProfileSetNotFlagged() {
	const (
		containerRuntimeDefault = "container-runtime-default"
		containerLocalhost      = "container-localhost"
		podRuntimeDefault       = "pod-runtime-default"
		podLocalhost            = "pod-localhost"
		containerOverridesPod   = "container-overrides-pod-unconfined"
	)

	s.addDeployment(containerRuntimeDefault, nil,
		[]v1.Container{containerWithProfile(runtimeDefaultProfile)}, nil, nil)
	localhostContainer := containerWithProfile(localhostProfile)
	localhostContainer.SecurityContext.SeccompProfile.LocalhostProfile = &localhostProfileRef
	s.addDeployment(containerLocalhost, nil,
		[]v1.Container{localhostContainer}, nil, nil)
	s.addDeployment(podRuntimeDefault,
		&v1.PodSecurityContext{SeccompProfile: &v1.SeccompProfile{Type: v1.SeccompProfileTypeRuntimeDefault}},
		[]v1.Container{{Name: containerName}}, nil, nil)
	s.addDeployment(podLocalhost,
		&v1.PodSecurityContext{SeccompProfile: &v1.SeccompProfile{Type: v1.SeccompProfileTypeLocalhost, LocalhostProfile: &localhostProfileRef}},
		[]v1.Container{{Name: containerName}}, nil, nil)
	// A container-level profile takes precedence over an "Unconfined" pod-level profile.
	s.addDeployment(containerOverridesPod,
		&v1.PodSecurityContext{SeccompProfile: &v1.SeccompProfile{Type: v1.SeccompProfileTypeUnconfined}},
		[]v1.Container{containerWithProfile(runtimeDefaultProfile)}, nil, nil)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param:       params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{},
		},
	})
}

func (s *SeccompProfileTestSuite) TestInitAndEphemeralContainers() {
	const (
		initNoProfile      = "init-container-no-profile"
		ephemeralNoProfile = "ephemeral-container-no-profile"
		podProfileCovers   = "pod-profile-covers-init-and-ephemeral"
	)

	compliant := containerWithProfile(runtimeDefaultProfile)
	s.addDeployment(initNoProfile, nil,
		[]v1.Container{compliant},
		[]v1.Container{{Name: initContainerName}}, nil)
	s.addDeployment(ephemeralNoProfile, nil,
		[]v1.Container{compliant}, nil,
		[]v1.EphemeralContainer{{EphemeralContainerCommon: v1.EphemeralContainerCommon{Name: ephemeralContainerName}}})
	// A pod-level profile covers init and ephemeral containers too.
	s.addDeployment(podProfileCovers,
		&v1.PodSecurityContext{SeccompProfile: &v1.SeccompProfile{Type: v1.SeccompProfileTypeRuntimeDefault}},
		[]v1.Container{{Name: containerName}},
		[]v1.Container{{Name: initContainerName}},
		[]v1.EphemeralContainer{{EphemeralContainerCommon: v1.EphemeralContainerCommon{Name: ephemeralContainerName}}})

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				initNoProfile:      expectedDiag(initContainerName),
				ephemeralNoProfile: expectedDiag(ephemeralContainerName),
			},
		},
	})
}
