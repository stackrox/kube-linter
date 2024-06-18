package volumeclaimtemplates

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/volumeclaimtemplates/internal/params"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
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
		Parameters: params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				k8sObj, ok := object.K8sObject.(*unstructured.Unstructured)
				if !ok {
					return nil
				}
				if k8sObj.GetKind() != "StatefulSet" {
					return nil
				}
				var statefulSet v1.StatefulSet
				err := runtime.DefaultUnstructuredConverter.FromUnstructured(k8sObj.UnstructuredContent(), &statefulSet)
				if err != nil {
					return nil
				}
				for _, volumeClaimTemplate := range statefulSet.Spec.VolumeClaimTemplates {
					if annotationValue, found := volumeClaimTemplate.Annotations[p.Annotation]; found {
						return []diagnostic.Diagnostic{{
							Message: fmt.Sprintf("found annotation %q with value %q in VolumeClaimTemplate", p.Annotation, annotationValue),
						}}
					}
				}
				return nil
			}, nil
	 } ),
	})
}
