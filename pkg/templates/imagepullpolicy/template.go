package livenessprobe

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/imagepullpolicy/internal/params"
	"golang.stackrox.io/kube-linter/pkg/templates/util"
	v1 "k8s.io/api/core/v1"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Image Pull Policy",
		Key:         "image-pull-policy",
		Description: "Flag containers with forbidden image pull policy",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			return util.PerContainerCheck(func(container *v1.Container) []diagnostic.Diagnostic {
				forbiddenPolicies := map[string]bool{}
				for _, v := range p.ForbiddenPolicies {
					forbiddenPolicies[v] = true
				}
				if _, ok := forbiddenPolicies[string(container.ImagePullPolicy)]; ok {
					return []diagnostic.Diagnostic{{Message: fmt.Sprintf("container %q has a %s pull image policy", container.Name, container.ImagePullPolicy)}}
				}
				return nil
			}), nil
		}),
	})
}
