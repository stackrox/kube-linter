package containercapabilities

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/diagnostic"
	"golang.stackrox.io/kube-linter/internal/matcher"
	"golang.stackrox.io/kube-linter/internal/objectkinds"
	"golang.stackrox.io/kube-linter/internal/templates"
	"golang.stackrox.io/kube-linter/internal/templates/containercapabilities/internal/params"
	"golang.stackrox.io/kube-linter/internal/templates/util"
	v1 "k8s.io/api/core/v1"
)

func checkForbiddenAdds(containerName string, paramCaps []string, scCaps []v1.Capability) (*diagnostic.Diagnostic, error) {
	for _, forbiddenCap := range paramCaps {
		forbiddenCapMatcher, err := matcher.ForString(forbiddenCap)
		if err != nil {
			return nil, err
		}
		for _, scCap := range scCaps {
			if forbiddenCapMatcher(string(scCap)) {
				// If any scCap exist in paramCaps, flag the container
				return &diagnostic.Diagnostic{
					Message: fmt.Sprintf("container %q has ADD capabilities: %q, which violates"+
						" the forbidden ADD capabilities for containers: %q",
						containerName,
						scCaps,
						paramCaps),
				}, nil
			}
		}
	}

	return nil, nil
}

func checkRequiredDrops(containerName string, paramCaps []string, scCaps []v1.Capability) (*diagnostic.Diagnostic, error) {
	for _, requiredCap := range paramCaps {
		requiredCapMatcher, err := matcher.ForString(requiredCap)
		if err != nil {
			return nil, err
		}
		found := false
		for _, scCap := range scCaps {
			if requiredCapMatcher(string(scCap)) {
				// Found this required cap, go check the next one
				found = true
				break
			}
		}
		if !found {
			// If any required drops do not exist in scCaps, flag the container
			return &diagnostic.Diagnostic{
				Message: fmt.Sprintf("container %q has DROP capabilities: %q, which does not "+
					"satisfy the required DROP capabilities for containers: %q",
					containerName,
					scCaps,
					paramCaps),
			}, nil
		}
	}

	return nil, nil
}

func init() {
	templates.Register(check.Template{
		HumanName:   "Verify container capabilities",
		Key:         "verify-container-capabilities",
		Description: "Flag containers that do not match capabilities requirements",
		SupportedObjectKinds: check.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			var returnedErr error
			return util.PerContainerCheck(func(container *v1.Container) []diagnostic.Diagnostic {
				// 2 = 1 for forbiddenAdds + 1 for allRequiredDrops
				result := make([]diagnostic.Diagnostic, 0, 2)
				sc := container.SecurityContext
				if sc != nil && sc.Capabilities != nil {
					diagnostic, err := checkForbiddenAdds(container.Name, p.ForbiddenAdds, sc.Capabilities.Add)
					if err != nil {
						returnedErr = errors.Wrap(err, "checking forbidden ADD capabilities")
						return nil
					}
					if diagnostic != nil {
						result = append(result, *diagnostic)
					}
					diagnostic, err = checkRequiredDrops(container.Name, p.RequiredDrops, sc.Capabilities.Drop)
					if err != nil {
						returnedErr = errors.Wrap(err, "checking required DROP capabilities")
						return nil
					}
					if diagnostic != nil {
						result = append(result, *diagnostic)
					}
				}
				return result
			}), returnedErr
		}),
	})
}
