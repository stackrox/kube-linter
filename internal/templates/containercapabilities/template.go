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

const (
	reservedCapabilitiesAll             = "all"
	matchLiteralReservedCapabilitiesAll = "(?i)" + reservedCapabilitiesAll
)

var (
	literalReservedCapabilitiesAllMatcher func(string) bool
)

func checkCapabilityDropList(
	containerName string,
	paramCapMatchers map[string]func(string) bool,
	scCaps *v1.Capabilities,
	forbidAll bool,
) *diagnostic.Diagnostic {
	if forbidAll {
		// Specified to drop all. Instead of trying to find a match to this paramCap, we should
		// make sure the special ReservedCapabilitiesAll is found in the DROP list as well
		for _, scCap := range scCaps.Drop {
			if literalReservedCapabilitiesAllMatcher(string(scCap)) {
				// Found "all" in drop list as well
				return nil
			}
		}
		return &diagnostic.Diagnostic{
			Message: fmt.Sprintf(
				"container %q has DROP capabilities: %q, but in fact all capabilities are required to be dropeed",
				containerName,
				scCaps.Drop),
		}
	}

	// Every forbidden capability specified by param should be found in the DROP list
	for paramCap, paramCapMatcher := range paramCapMatchers {
		var found bool
		for _, scCap := range scCaps.Drop {
			// User can specify to drop "all" under containers as well
			if paramCapMatcher(string(scCap)) || literalReservedCapabilitiesAllMatcher(string(scCap)) {
				// This forbidden cap exists in the DROP list. Check for the next forbidden cap
				found = true
				break
			}
		}
		if !found {
			return &diagnostic.Diagnostic{
				Message: fmt.Sprintf("container %q has DROP capabilities: %q, but does not drop "+
					"capability %q which is required",
					containerName,
					scCaps.Drop,
					paramCap),
			}
		}
	}
	return nil
}

func checkCapabilityAddList(
	containerName string,
	paramCapMatchers map[string]func(string) bool,
	scCaps *v1.Capabilities,
	forbidAll bool,
	exceptionCapMatchers map[string]func(string) bool,
) *diagnostic.Diagnostic {
	if forbidAll {
		// User has forbidden all capabilities
		for _, scCap := range scCaps.Add {
			var excluded bool
			for _, exceptionCapMatcher := range exceptionCapMatchers {
				if exceptionCapMatcher(string(scCap)) {
					// Forgive this capability
					excluded = true
					break
				}
			}
			if !excluded {
				return &diagnostic.Diagnostic{
					Message: fmt.Sprintf(
						"container %q has ADD capability: %q, but no capabilities should be added at all and"+
							" this capabilty is not included in the exceptions list",
						containerName,
						scCap),
				}
			}
		}
		// No violations
		return nil
	}

	// Any capability from scCaps should not match with any from paramCaps
	for _, paramCapMatcher := range paramCapMatchers {
		for _, scCap := range scCaps.Add {
			// User can specify to add "all" under containers as well.
			if paramCapMatcher(string(scCap)) || literalReservedCapabilitiesAllMatcher(string(scCap)) {
				// A capability from ADD list matched with a cap from forbidden capabilities list.
				return &diagnostic.Diagnostic{
					Message: fmt.Sprintf("container %q has ADD capability: %q, which matched with "+
						"the forbidden capability for containers",
						containerName,
						scCap),
				}
			}
		}
	}
	return nil
}

func checkForbidAll(paramCaps []string) bool {
	for _, paramCap := range paramCaps {
		if literalReservedCapabilitiesAllMatcher(paramCap) {
			return true
		}
	}
	return false
}

func verifyExceptionsList(forbidAll bool, exceptions []string) error {
	// Check if forbidAll is set
	if !forbidAll && len(exceptions) != 0 {
		return errors.Errorf("for verifying container capabilities, \"Exceptions\" list should only"+
			" be filled when %q capabilities specified in the forbidden list", reservedCapabilitiesAll)
	}
	// Check no "all" in exceptions list
	for _, cap := range exceptions {
		if literalReservedCapabilitiesAllMatcher(cap) {
			return errors.Errorf("capabilities exceptions list should not contain %q", reservedCapabilitiesAll)
		}
	}
	return nil
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
			var err error
			literalReservedCapabilitiesAllMatcher, err = matcher.ForString(matchLiteralReservedCapabilitiesAll)
			if err != nil {
				return nil, err
			}

			// Check for "all" in case user does not want any capability in containers.
			forbidAll := checkForbidAll(p.ForbiddenCapabilities)
			err = verifyExceptionsList(forbidAll, p.Exceptions)
			if err != nil {
				return nil, err
			}

			paramCapMatchers := make(map[string]func(string) bool, len(p.ForbiddenCapabilities))
			exceptionCapMatchers := make(map[string]func(string) bool, len(p.Exceptions))
			if !forbidAll {
				for _, cap := range p.ForbiddenCapabilities {
					capMatcher, err := matcher.ForString(cap)
					if err != nil {
						return nil, errors.Wrapf(err, "checking container capabilities. invalid capability: %s", cap)
					}
					paramCapMatchers[cap] = capMatcher
				}
			} else {
				for _, cap := range p.Exceptions {
					capMatcher, err := matcher.ForString(cap)
					if err != nil {
						return nil, errors.Wrapf(err, "checking container capabilities. invalid capability: %s", cap)
					}
					exceptionCapMatchers[cap] = capMatcher
				}
			}

			return util.PerContainerCheck(func(container *v1.Container) []diagnostic.Diagnostic {
				// 2 = 1 for forbiddenAdds + 1 for allRequiredDrops
				result := make([]diagnostic.Diagnostic, 0, 2)
				sc := container.SecurityContext
				if sc != nil && sc.Capabilities != nil {
					// Check none of the forbidden capabilities exists in ADDs
					diagnostic :=
						checkCapabilityAddList(
							container.Name,
							paramCapMatchers,
							sc.Capabilities,
							forbidAll,
							exceptionCapMatchers)
					if diagnostic != nil {
						result = append(result, *diagnostic)
					}
					// Check every forbidden capabilities should exist in DROPs
					diagnostic =
						checkCapabilityDropList(container.Name, paramCapMatchers, sc.Capabilities, forbidAll)
					if diagnostic != nil {
						result = append(result, *diagnostic)
					}
				}
				return result
			}), nil
		}),
	})
}
