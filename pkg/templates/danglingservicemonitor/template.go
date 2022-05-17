package danglingservicemonitor

import (
	"fmt"

	k8sMonitoring "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/danglingservicemonitor/internal/params"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Dangling Service Monitor",
		Key:         "dangling-servicemonitor",
		Description: "Flag service monitors which do not match any service",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.ServiceMonitor},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			return func(lintCtx lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				serviceMonitor, ok := object.K8sObject.(*k8sMonitoring.ServiceMonitor)
				if !ok {
					return nil
				}
				nsSelector := serviceMonitor.Spec.NamespaceSelector
				nsSelectorSet := len(nsSelector.MatchNames) != 0 || nsSelector.Any

				labelSelectors := serviceMonitor.Spec.Selector.MatchLabels
				labelSelectorSet := len(labelSelectors) != 0
				if !labelSelectorSet && !nsSelectorSet {
					return []diagnostic.Diagnostic{{
						Message: "service monitor has no selector specified",
					}}
				}
				labelSelector, err := metaV1.LabelSelectorAsSelector(&metaV1.LabelSelector{MatchLabels: serviceMonitor.Spec.Selector.MatchLabels})
				if err != nil {
					return []diagnostic.Diagnostic{{
						Message: fmt.Sprintf("service monitor has invalid label selector: %v", err),
					}}
				}
				for _, obj := range lintCtx.Objects() {
					services, found := obj.K8sObject.(*v1.Service)
					if !found {
						continue
					}
					if checkNamespaceSelector(nsSelector, services) {
						if !labelSelectorSet {
							return nil
						}
						if labelSelectorSet && labelSelector.Matches(labels.Set(services.Labels)) {
							return nil
						} else {
							continue
						}
					}
					if labelSelector.Matches(labels.Set(services.Labels)) && labelSelectorSet && !nsSelectorSet {
						// Found!
						return nil
					}

				}
				return []diagnostic.Diagnostic{{Message: fmt.Sprintf("no services found matching the service monitor's label selector (%s) and namespace selector (%s)", labelSelector, nsSelector.MatchNames)}}
			}, nil
		}),
	})
}

func checkNamespaceSelector(namespaceSelector k8sMonitoring.NamespaceSelector, service *v1.Service) bool {
	if namespaceSelector.Any {
		return true
	}
	for _, ns := range namespaceSelector.MatchNames {
		if ns == service.Namespace {
			return true
		}
	}
	return false
}
