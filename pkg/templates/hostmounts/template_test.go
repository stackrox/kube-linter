package hostmounts

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/hostmounts/internal/params"
	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

const (
	dirUsr     = "/usr"
	dirEtc     = "/etc"
	container0 = "container0"
)

var dirParams = []string{"^/usr$", "^/etc$"}

func TestHostMounts(t *testing.T) {
	suite.Run(t, new(HostMountsTestSuite))
}

type HostMountsTestSuite struct {
	templates.TemplateTestSuite
	ctx *mocks.MockLintContext
}

func (s *HostMountsTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *HostMountsTestSuite) addDeploymentWithContainer(name, containername string) {
	s.ctx.AddMockDeployment(s.T(), name)
	s.ctx.ModifyDeployment(s.T(), name, func(deployment *appsV1.Deployment) {
		c := v1.Container{
			Name:            containername,
			Image:           "myorg/myimage:tag",
			Resources:       v1.ResourceRequirements{},
			VolumeMounts:    []v1.VolumeMount{},
			SecurityContext: &v1.SecurityContext{},
		}
		deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, c)
	})
}

func (s *HostMountsTestSuite) TestWithoutVolumes() {
	s.addDeploymentWithContainer("no-volume-dep", container0)
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				Dirs: []string{},
			},
			Diagnostics:              nil,
			ExpectInstantiationError: false,
		},
	})
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				Dirs: dirParams,
			},
			Diagnostics:              nil,
			ExpectInstantiationError: false,
		},
	})
}

func (s *HostMountsTestSuite) addDeploymentWithVolume(name, vname, vpath, mpath string) {
	s.addDeploymentWithContainer(name, container0)
	s.ctx.ModifyDeployment(s.T(), name, func(deployment *appsV1.Deployment) {
		deployment.Spec.Template.Labels = map[string]string{"app": name}
		vol := v1.Volume{
			Name: vname,
			VolumeSource: v1.VolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: vpath,
				},
			},
		}
		deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, vol)
		mount := v1.VolumeMount{
			Name:      vname,
			ReadOnly:  false,
			MountPath: mpath,
		}
		deployment.Spec.Template.Spec.Containers[0].VolumeMounts = append(deployment.Spec.Template.Spec.Containers[0].VolumeMounts, mount)
	})
}

func (s *HostMountsTestSuite) TestWithVolumes() {
	const (
		oneNonSensitiveVolume = "one-non-sensitive-volume"
		oneSensitiveVolume    = "one-sensitive-volume"
		twoSensitiveVolumes   = "two-sensitive-volumes"
	)
	p := params.Params{
		Dirs: dirParams,
	}

	s.addDeploymentWithVolume(oneNonSensitiveVolume, "one-non-sensitive-volume", "/mydata", "non-sensitive-mount")
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param:                    p,
			Diagnostics:              nil,
			ExpectInstantiationError: false,
		},
	})

	s.addDeploymentWithVolume(oneSensitiveVolume, "one-sensitive-volume", dirUsr, "one-sensitive-mount")
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: p,
			Diagnostics: map[string][]diagnostic.Diagnostic{
				oneSensitiveVolume: {
					{Message: fmt.Sprintf("host system directory %q is mounted on container %q", dirUsr, container0)},
				},
			},
			ExpectInstantiationError: false,
		},
	})
	s.addDeploymentWithVolume(twoSensitiveVolumes, "two-sensitive-volumes", dirEtc, "two-sensitive-mounts")
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: p,
			Diagnostics: map[string][]diagnostic.Diagnostic{
				oneSensitiveVolume: {
					{Message: fmt.Sprintf("host system directory %q is mounted on container %q", dirUsr, container0)},
				},
				twoSensitiveVolumes: {
					{Message: fmt.Sprintf("host system directory %q is mounted on container %q", dirEtc, container0)},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}
