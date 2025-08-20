package volumeclaimtemplates

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/volumeclaimtemplates/internal/params"
)

const (
	templateKey = "statefulset-volumeclaimtemplate-annotation"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "StatefulSet VolumeClaimTemplate Annotation",
		Key:         templateKey,
		Description: "Check if StatefulSet's VolumeClaimTemplate contains a specific annotation",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				sts, ok := extract.StatefulSetSpec(object.K8sObject)
				if !ok {
					return nil
				}
				var diagnostics []diagnostic.Diagnostic
				for _, vct := range sts.VolumeClaimTemplates {
					if vct.Annotations == nil || vct.Annotations[p.Annotation] == "" {
						diagnostics = append(diagnostics, diagnostic.Diagnostic{
							Message: fmt.Sprintf("StatefulSet's VolumeClaimTemplate is missing required annotation: %s", p.Annotation),
						})
					}
				}
				return diagnostics
			}, nil
		}),
	})
}
