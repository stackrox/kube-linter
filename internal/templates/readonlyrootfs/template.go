package readonlyrootfs

import (
	"fmt"

	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/diagnostic"
	"golang.stackrox.io/kube-linter/internal/extract"
	"golang.stackrox.io/kube-linter/internal/lintcontext"
	"golang.stackrox.io/kube-linter/internal/objectkinds"
	"golang.stackrox.io/kube-linter/internal/templates"
	"golang.stackrox.io/kube-linter/internal/templates/readonlyrootfs/internal/params"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Read-only Root Filesystems",
		Key:         "read-only-root-fs",
		Description: "Flag containers without read-only root file systems",
		SupportedObjectKinds: check.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			return func(_ *lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				podSpec, found := extract.PodSpec(object.K8sObject)
				if !found {
					return nil
				}
				var results []diagnostic.Diagnostic
				for _, container := range podSpec.Containers {
					sc := container.SecurityContext
					if sc == nil || sc.ReadOnlyRootFilesystem == nil || !*sc.ReadOnlyRootFilesystem {
						results = append(results, diagnostic.Diagnostic{Message: fmt.Sprintf("container %q does not have a read-only root file system", container.Name)})
					}
				}
				return results
			}, nil
		}),
	})
}
