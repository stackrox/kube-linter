package dnsconfigoptions

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/dnsconfigoptions/internal/params"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "DnsConfig Options",
		Key:         "dnsconfig-options",
		Description: "Flag objects that don't have specified DNSConfig Options",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				podTemplateSpec, found := extract.PodTemplateSpec(object.K8sObject)
				if !found {
					return nil
				}
				if podTemplateSpec.Spec.DNSConfig == nil {
					return []diagnostic.Diagnostic{{Message: "Object does not define any DNSConfig rules."}}
				}
				if podTemplateSpec.Spec.DNSConfig.Options == nil {
					return []diagnostic.Diagnostic{{Message: "Object does not define any DNSConfig Options."}}
				}

				for _, option := range podTemplateSpec.Spec.DNSConfig.Options {
					if option.Name == p.Key && *option.Value == p.Value {
						// Found
						return nil
					}
				}
				return []diagnostic.Diagnostic{{
					Message: fmt.Sprintf("DNSConfig Options \"%s:%s\" not found.", p.Key, p.Value),
				}}
			}, nil
		}),
	})
}
