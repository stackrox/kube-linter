package run

import (
	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/errorhelpers"
	"golang.stackrox.io/kube-linter/internal/objectkinds"
	"golang.stackrox.io/kube-linter/internal/templates"
)

type instantiatedCheck struct {
	Name    string
	Func    check.Func
	Matcher objectkinds.Matcher
}

func validateAndInstantiate(c *check.Check) (*instantiatedCheck, error) {
	validationErrs := errorhelpers.NewErrorList("validating check")
	if c.Name == "" {
		validationErrs.AddString("no name specified")
	}
	template, found := templates.Get(c.Template)
	if !found {
		validationErrs.AddStringf("template %q not found", c.Template)
		return nil, validationErrs.ToError()
	}

	supportedParams := make(map[string]struct{}, len(template.Parameters))
	for _, param := range template.Parameters {
		if param.Required {
			if _, found := c.Params[param.ParamName]; !found {
				validationErrs.AddStringf("required param %q not specified", param.ParamName)
			}
		}
		supportedParams[param.ParamName] = struct{}{}
	}
	for passedParam := range c.Params {
		if _, isSupported := supportedParams[passedParam]; !isSupported {
			validationErrs.AddStringf("unknown param %q passed", passedParam)
		}
	}
	if err := validationErrs.ToError(); err != nil {
		return nil, err
	}

	i := &instantiatedCheck{Name: c.Name}
	var objectKinds check.ObjectKindsDesc
	if c.Scope != nil {
		objectKinds = *c.Scope
	} else {
		objectKinds = template.SupportedObjectKinds
	}
	matcher, err := objectkinds.ConstructMatcher(objectKinds.ObjectKindNames...)
	if err != nil {
		return nil, err
	}
	i.Matcher = matcher
	checkFunc, err := template.Instantiate(c.Params)
	if err != nil {
		return nil, errors.Wrap(err, "instantiating check")
	}
	i.Func = checkFunc
	return i, nil
}
