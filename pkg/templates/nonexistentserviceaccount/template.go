package nonexistentserviceaccount

import (
	"fmt"

	"golang.stackrox.io/kube-linter/internal/set"
	"golang.stackrox.io/kube-linter/internal/stringutils"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/nonexistentserviceaccount/internal/params"
	v1 "k8s.io/api/core/v1"
)

var (
	serviceAccountGVK = v1.SchemeGroupVersion.WithKind("ServiceAccount")
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Non-Existent Service Account",
		Key:         "non-existent-service-account",
		Description: "Flag cases where a pod references a non-existent service account",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(_ params.Params) (check.Func, error) {
			return func(lintCtx lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				podSpec, found := extract.PodSpec(object.K8sObject)
				if !found {
					return nil
				}
				sa := stringutils.OrDefault(podSpec.ServiceAccountName, podSpec.DeprecatedServiceAccount)
				// Default service account always exists.
				if sa == "" || sa == "default" {
					return nil
				}
				ns := object.K8sObject.GetNamespace()
				serviceAccountsInCtx := set.NewStringSet()
				for _, otherObj := range lintCtx.Objects() {
					k8sObj := otherObj.K8sObject
					if k8sObj.GetObjectKind().GroupVersionKind() == serviceAccountGVK && k8sObj.GetNamespace() == ns {
						serviceAccountsInCtx.Add(k8sObj.GetName())
					}
				}
				if !serviceAccountsInCtx.Contains(sa) {
					return []diagnostic.Diagnostic{{Message: fmt.Sprintf(
						"serviceAccount %q not found", sa)}}
				}
				return nil
			}, nil
		}),
	})
}
