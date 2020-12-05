package check

import (
	"golang.stackrox.io/kube-linter/internal/pointers"
)

// ParameterType represents the expected type of a particular parameter.
type ParameterType string

// This block enumerates all known type names.
// These type names are chosen to be aligned with OpenAPI/JSON schema.
const (
	StringType  ParameterType = "string"
	IntegerType ParameterType = "integer"
	BooleanType ParameterType = "boolean"
	NumberType  ParameterType = "number"
	ObjectType  ParameterType = "object"
	ArrayType   ParameterType = "array"
)

// ParameterDesc describes a parameter.
type ParameterDesc struct {
	Name        string
	Type        ParameterType
	Description string

	Examples []string

	// Enum is set if the object is always going to be one of a specified set of values.
	// Only relevant if Type is "string"
	Enum []string

	// SubParameters are the child parameters of the given parameter.
	// Only relevant if Type is "object".
	SubParameters []ParameterDesc

	// ArrayElemType is only set when the object is of type array, and it describes the type
	// of the element of the array.
	ArrayElemType ParameterType

	// Required denotes whether the parameter is required.
	Required bool

	// NoRegex is set if the parameter does not support regexes.
	// Only relevant if Type is "string".
	NoRegex bool

	// NotNegatable is set if the parameter does not support negation via a leading !.
	// OnlyRelevant if Type is "string".
	NotNegatable bool

	// Fields below are for internal use only.

	XXXStructFieldName string
	XXXIsPointer       bool
}

// HumanReadableParamDesc is a human-friendly representation of a ParameterDesc.
// It is intended only for API documentation/JSON marshaling, and must NOT be used for
// any business logic.
type HumanReadableParamDesc struct {
	Name            string                   `json:"name"`
	Type            ParameterType            `json:"type"`
	Description     string                   `json:"description"`
	Required        bool                     `json:"required"`
	Examples        []string                 `json:"examples,omitempty"`
	RegexAllowed    *bool                    `json:"regexAllowed,omitempty"`
	NegationAllowed *bool                    `json:"negationAllowed,omitempty"`
	SubParameters   []HumanReadableParamDesc `json:"subParameters,omitempty"`
	ArrayElemType   ParameterType            `json:"arrayElemType,omitempty"`
}

// HumanReadableFields returns a human-friendly representation of this ParameterDesc.
func (p *ParameterDesc) HumanReadableFields() HumanReadableParamDesc {
	out := HumanReadableParamDesc{
		Name:        p.Name,
		Type:        p.Type,
		Description: p.Description,
		Required:    p.Required,
		Examples:    p.Examples,
	}

	if p.Type == StringType ||
		(p.Type == ArrayType && p.ArrayElemType == StringType) {
		out.RegexAllowed = pointers.Bool(!p.NoRegex)
		out.NegationAllowed = pointers.Bool(!p.NotNegatable)
	}

	if p.Type == ArrayType {
		out.ArrayElemType = p.ArrayElemType
	}

	if len(p.SubParameters) > 0 {
		subParamFields := make([]HumanReadableParamDesc, 0, len(p.SubParameters))
		for _, subParam := range p.SubParameters {
			subParamFields = append(subParamFields, subParam.HumanReadableFields())
		}
		out.SubParameters = subParamFields
	}

	return out
}
