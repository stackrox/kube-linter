package params

import (
	"fmt"
	"strings"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/templates/util"
)

// Params represents the params accepted by this template.
type Params struct {
	// Recursive determines whether to check keys recursively at all nesting levels.
	// Default is true.
	Recursive bool
}

var (
	recursiveParamDesc = check.ParameterDesc{
		Name:        "recursive",
		Type:        check.BooleanType,
		Description: "Check keys recursively at all nesting levels. Default is true.",
		Required:    false,
	}

	// ParamDescs is the list of parameter descriptors for this template.
	ParamDescs = []check.ParameterDesc{
		recursiveParamDesc,
	}
)

// Validate validates the parameters.
func (p *Params) Validate() error {
	var validationErrors []string
	if len(validationErrors) > 0 {
		return fmt.Errorf("invalid parameters: %s", strings.Join(validationErrors, ", "))
	}
	return nil
}

// ParseAndValidate instantiates a Params object out of the passed map[string]interface{},
// validates it, and returns it.
func ParseAndValidate(m map[string]interface{}) (interface{}, error) {
	var p Params
	// Set defaults
	p.Recursive = true

	if err := util.DecodeMapStructure(m, &p); err != nil {
		return nil, err
	}
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return p, nil
}

// WrapInstantiateFunc is a convenience wrapper that wraps an untyped instantiate function
// into a typed one.
func WrapInstantiateFunc(f func(p Params) (check.Func, error)) func(interface{}) (check.Func, error) {
	return func(paramsInt interface{}) (check.Func, error) {
		return f(paramsInt.(Params))
	}
}
