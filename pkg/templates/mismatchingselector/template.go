package mismatchingselector

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/mismatchingselector/internal/params"
	v1 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Mismatching Selector",
		Key:         "mismatching-selector",
		Description: "Flag deployments where the selector doesn't match the labels in the pod template spec",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(_ params.Params) (check.Func, error) {
			return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				selector, found := extract.Selector(object.K8sObject)
				if !found {
					return nil
				}
				if selector == nil || (len(selector.MatchLabels) == 0 && len(selector.MatchExpressions) == 0) {
					switch object.K8sObject.(type) {
					// It's okay for CronJobs and Jobs not to have selectors.
					case *v1beta1.CronJob, *v1.Job, *v1.CronJob:
						return nil
					}
					return []diagnostic.Diagnostic{{
						Message: "object has no selector specified",
					}}
				}

				podTemplateSpec, hasPods := extract.PodTemplateSpec(object.K8sObject)
				if !hasPods {
					return nil
				}
				labelSelector, err := metaV1.LabelSelectorAsSelector(selector)
				if err != nil {
					return []diagnostic.Diagnostic{{
						Message: fmt.Sprintf("object has invalid label selector: %v", err),
					}}
				}
				if labelSelector.Matches(labels.Set(podTemplateSpec.Labels)) {
					return nil
				}
				return []diagnostic.Diagnostic{{Message: fmt.Sprintf("labels in pod spec (%v) do not match labels in selector (%v)", podTemplateSpec.Labels, selector)}}
			}, nil
		}),
	})
}
