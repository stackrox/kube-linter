package targetport

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/extract/customtypes"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/targetport/internal/params"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	k8sValidation "k8s.io/apimachinery/pkg/util/validation"
)

const (
	templateKey = "target-port"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Target Port",
		Key:         templateKey,
		Description: "Flag containers and services using not allowed port names or numbers",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike, objectkinds.Service},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(_ params.Params) (check.Func, error) {
			return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				podSpec, foundPodSpec := extract.PodSpec(object.K8sObject)
				if foundPodSpec {
					return findPodPorts(&podSpec)
				}

				service, foundService := object.K8sObject.(*coreV1.Service)
				if foundService {
					return findServicePorts(service)
				}

				return nil
			}, nil

		}),
	})
}

func findPodPorts(podSpec *customtypes.PodSpec) []diagnostic.Diagnostic {
	var results []diagnostic.Diagnostic

	containers := podSpec.AllContainers()
	for _, container := range containers {
		for _, port := range container.Ports {
			if port.Name == "" {
				continue
			}

			violations := k8sValidation.IsValidPortName(port.Name)
			for _, violation := range violations {
				results = append(results, diagnostic.Diagnostic{
					Message: fmt.Sprintf("port name %q in container %q %s",
						port.Name, container.Name, violation),
				})
			}
		}
	}

	return results
}

func findServicePorts(service *coreV1.Service) []diagnostic.Diagnostic {
	var results []diagnostic.Diagnostic

	for _, port := range service.Spec.Ports {
		targetPort := port.TargetPort
		if targetPort.Type == intstr.Int && targetPort.IntVal == 0 {
			continue
		}

		if targetPort.Type == intstr.String {
			violations := k8sValidation.IsValidPortName(targetPort.String())
			for _, violation := range violations {
				results = append(results, diagnostic.Diagnostic{
					Message: fmt.Sprintf("port targetPort %q in service %q %s",
						targetPort.String(), service.Name, violation),
				})
			}
		}

		if targetPort.Type == intstr.Int {
			violations := k8sValidation.IsValidPortNum(targetPort.IntValue())
			for _, violation := range violations {
				results = append(results, diagnostic.Diagnostic{
					Message: fmt.Sprintf("port targetPort %q in service %q %s",
						targetPort.String(), service.Name, violation),
				})
			}
		}
	}

	return results
}
