package readonlyrootfs

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/readonlyrootfs/internal/params"
	"golang.stackrox.io/kube-linter/pkg/templates/util"
	v1 "k8s.io/api/core/v1"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Read-only Root Filesystems",
		Key:         "read-only-root-fs",
		Description: "Flag containers without read-only root file systems",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			return util.PerContainerCheck(func(container *v1.Container) []diagnostic.Diagnostic {
				sc := container.SecurityContext
				if sc == nil || sc.ReadOnlyRootFilesystem == nil || !*sc.ReadOnlyRootFilesystem {
					return []diagnostic.Diagnostic{{Message: fmt.Sprintf("container %q does not have a read-only root file system", container.Name)}}
				}
				return nil
			}), nil
		}),
	})
}
