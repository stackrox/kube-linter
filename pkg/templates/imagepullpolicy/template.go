package imagepullpolicy

import (
	"fmt"

	"golang.stackrox.io/kube-linter/internal/set"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/imagepullpolicy/internal/params"
	"golang.stackrox.io/kube-linter/pkg/templates/util"
	v1 "k8s.io/api/core/v1"
)

const (
	templateKey = "image-pull-policy"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Image Pull Policy",
		Key:         templateKey,
		Description: "Flag containers with forbidden image pull policy",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			forbiddenPolicies := set.NewStringSet(p.ForbiddenPolicies...)
			return util.PerContainerCheck(func(container *v1.Container) []diagnostic.Diagnostic {
				if forbiddenPolicies.Contains(string(container.ImagePullPolicy)) {
					return []diagnostic.Diagnostic{{Message: fmt.Sprintf("container %q has imagePullPolicy set to %s", container.Name, container.ImagePullPolicy)}}
				}
				return nil
			}), nil
		}),
	})
}
