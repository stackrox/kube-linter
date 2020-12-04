package containercapabilities

import (
	"fmt"
	"strings"

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

const containerCapabilitiesDelim = ","

type capabilityCheckFunc func(capMatcher func(string) bool, scCap v1.Capability) bool

// Takes in a user provided checkFunc to determine if this function should return or not.
//   - If checkFunc returns false, this function returns immediately with value false
//   - If checkFunc returns true, this function keeps on looping until it hits a false,
//     or it finishes traversing all capabilities and returns true to indicate check has passed
func checkCapabilities(caps []string, scCaps []v1.Capability, checkFunc capabilityCheckFunc) (bool, error) {
	for _, cap := range caps {
		capMatcher, err := matcher.ForString(cap)
		if err != nil {
			return false, errors.Wrapf(err, "invalid capability: %s", cap)
		}
		for _, scCap := range scCaps {
			if !checkFunc(capMatcher, scCap) {
				return false, nil
			}
		}
	}
	// Traversed through all of caps, return true to indicate check has passed
	return true, nil
}

func checkForbiddenAdds(containerName string, paramCaps []string, scCaps []v1.Capability) (*diagnostic.Diagnostic, error) {
	passed, err :=
		checkCapabilities(paramCaps, scCaps, func(capMatcher func(string) bool, scCap v1.Capability) bool {
			// If any matches, then the check should fail
			//Otherwise keep checking
			return !capMatcher(string(scCap))
		})
	if err != nil {
		return nil, err
	}
	if !passed {
		return &diagnostic.Diagnostic{
			Message: fmt.Sprintf("container %q has ADD capabilities: %q, which violates"+
				" the forbidden ADD capabilities for containers: %q",
				containerName,
				scCaps,
				paramCaps),
		}, nil
	}
	return nil, nil
}

func checkRequiredDrops(containerName string, paramCaps []string, scCaps []v1.Capability) (*diagnostic.Diagnostic, error) {
	passed, err :=
		checkCapabilities(paramCaps, scCaps, func(capMatcher func(string) bool, scCap v1.Capability) bool {
			// If any mismatches, then the check should fail
			// Otherwise keep checking
			return capMatcher(string(scCap))
		})
	if err != nil {
		return nil, err
	}
	if !passed {
		return &diagnostic.Diagnostic{
			Message: fmt.Sprintf("container %q has DROP capabilities: %q, which does not "+
				"satisfy the required DROP capabilities for containers: %q",
				containerName,
				scCaps,
				paramCaps),
		}, nil
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
			allForbiddenAdds := strings.Split(p.ForbiddenAdds, containerCapabilitiesDelim)
			allRequiredDrops := strings.Split(p.RequiredDrops, containerCapabilitiesDelim)
			var returnedErr error
			return util.PerContainerCheck(func(container *v1.Container) []diagnostic.Diagnostic {
				// 2 = 1 for forbiddenAdds + 1 for allRequiredDrops
				result := make([]diagnostic.Diagnostic, 0, 2)
				sc := container.SecurityContext
				if sc != nil && sc.Capabilities != nil {
					diagnostic, err := checkForbiddenAdds(container.Name, allForbiddenAdds, sc.Capabilities.Add)
					if err != nil {
						returnedErr = errors.Wrap(err, "checking forbidden ADD capabilities")
						return nil
					}
					if diagnostic != nil {
						result = append(result, *diagnostic)
					}
					diagnostic, err = checkRequiredDrops(container.Name, allRequiredDrops, sc.Capabilities.Drop)
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
