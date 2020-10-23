package instantiatedcheck

import (
	"regexp"

	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/errorhelpers"
	"golang.stackrox.io/kube-linter/internal/objectkinds"
	"golang.stackrox.io/kube-linter/internal/templates"
)

// An InstantiatedCheck is the runtime instantiation of a check, which fuses the metadata in a check
// spec with the runtime information from a template.
type InstantiatedCheck struct {
	Name    string
	Func    check.Func
	Matcher objectkinds.Matcher
}

var (
	validCheckNameRegex = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
)

// ValidateAndInstantiate validates the check, and creates an instantiated check if the check
// is valid.
func ValidateAndInstantiate(c *check.Check) (*InstantiatedCheck, error) {
	validationErrs := errorhelpers.NewErrorList("validating check")
	if c.Name == "" {
		validationErrs.AddString("no name specified")
	}
	if !validCheckNameRegex.MatchString(c.Name) {
		validationErrs.AddStringf("invalid name %s, must match regex %s", c.Name, validCheckNameRegex.String())
	}
	template, found := templates.Get(c.Template)
	if !found {
		validationErrs.AddStringf("template %q not found", c.Template)
		return nil, validationErrs.ToError()
	}

	params, err := template.ParseAndValidateParams(c.Params)
	if err != nil {
		return nil, errors.Wrap(err, "validating and instantiating params")
	}

	if err := validationErrs.ToError(); err != nil {
		return nil, err
	}

	i := &InstantiatedCheck{Name: c.Name}
	var objectKinds check.ObjectKindsDesc
	if c.Scope != nil {
		objectKinds = *c.Scope
	} else {
		objectKinds = template.SupportedObjectKinds
	}
	matcher, err := objectkinds.ConstructMatcher(objectKinds.ObjectKinds...)
	if err != nil {
		return nil, err
	}
	i.Matcher = matcher
	checkFunc, err := template.Instantiate(params)
	if err != nil {
		return nil, errors.Wrap(err, "instantiating check")
	}
	i.Func = checkFunc
	return i, nil
}
