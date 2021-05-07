package hostipc

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/hostipc/internal/params"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Host IPC",
		Key:         "host-ipc",
		Description: "Flag Pod sharing host's IPC namespace",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(_ params.Params) (check.Func, error) {
			return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				podSpec, found := extract.PodSpec(object.K8sObject)
				if !found {
					return nil
				}
				object.K8sObject.GetName()
				if podSpec.HostIPC {
					return []diagnostic.Diagnostic{{Message: fmt.Sprintf("resource %s shares,  host's IPC namespace.", object.K8sObject.GetName())}}
				}
				return nil
			}, nil
		}),
	})
}
