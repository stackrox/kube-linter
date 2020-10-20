package check

// ParameterType represents the expected type of a particular parameter.
type ParameterType string

// This block enumerates all known type names.
// These type names are chosen to be aligned with OpenAPI/JSON schema.
const (
	StringType  ParameterType = "string"
	IntegerType               = "integer"
	BooleanType               = "boolean"
	NumberType                = "number"
	ObjectType                = "object"
)

// ParameterDesc describes a parameter.
type ParameterDesc struct {
	Name        string
	Type        ParameterType
	Description string

	Examples []string

	// SubParameters are the child parameters of the given parameter.
	// Only relevant if Type is "object".
	SubParameters []ParameterDesc

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
}

// HumanReadableFields returns a human-friendly representation of this ParameterDesc
func (p *ParameterDesc) HumanReadableFields() map[string]interface{} {
	m := map[string]interface{}{
		"name": p.Name,
		"type": p.Type,
		"description": p.Description,
		"required": p.Required,
	}

	if len(p.Examples) > 0 {
		m["examples"] = p.Examples
	}

	if p.Type == StringType {
		m["regexAllowed"] = !p.NoRegex
		m["negationAllowed"] = !p.NotNegatable
	}

	if len(p.SubParameters) > 0 {
		subParamFields := make([]map[string]interface{}, 0, len(p.SubParameters))
		for _, subParam := range p.SubParameters {
			subParamFields = append(subParamFields, subParam.HumanReadableFields())
		}
		m["subParameters"] = subParamFields
	}
	return m
}