package deprecatedfield

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/deprecatedfield/internal/params"
)

const (
	templateKey         = "deprecated-field"
	deprecatedStructTag = "deprecated"
)

// findDeprecatedFields recursively finds all deprecated fields that are set in the object.
// It returns a list of field paths (e.g., "spec.admissionControl.listenOnCreates") which were
// specified even though they are deprecated.
func findDeprecatedFields(objValue reflect.Value, path string) []string {
	var deprecatedPaths []string

	// Dereference pointers.
	for objValue.Kind() == reflect.Pointer {
		if objValue.IsNil() {
			return nil
		}
		objValue = objValue.Elem()
	}

	// Only process structs.
	if !valueIsStruct(objValue) {
		return nil
	}

	// Iterate over all the struct fields.
	objType := objValue.Type()
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		fieldValue := objValue.Field(i)

		// Skip unexported fields.
		if !field.IsExported() {
			continue
		}

		jsonName := determineJsonNameForField(field)
		if jsonName == "" {
			continue
		}

		// Build the field path.
		var fieldPath string
		if path == "" {
			fieldPath = jsonName
		} else {
			fieldPath = path + "." + jsonName
		}

		isSet := !valueIsZero(fieldValue)
		isDeprecated := fieldIsDeprecated(field)
		isStruct := valueIsStruct(fieldValue)

		if isDeprecated && isSet {
			deprecatedPaths = append(deprecatedPaths, fieldPath)
			if isStruct {
				// We are dealing with a deprecated sub struct, it is sufficient to flag this once.
				// To not flag every additional field set in this sub struct, we don't recurse further.
				continue
			}
		}

		// Recursively process nested structs.
		nestedPaths := findDeprecatedFields(fieldValue, fieldPath)
		deprecatedPaths = append(deprecatedPaths, nestedPaths...)
	}

	return deprecatedPaths
}

// fieldIsDeprecated checks if a struct field is marked has "deprecated" struct tag set to "true".
func fieldIsDeprecated(field reflect.StructField) bool {
	deprecatedTagVal, _ := field.Tag.Lookup(deprecatedStructTag)
	isDeprecated, _ := strconv.ParseBool(deprecatedTagVal)
	return isDeprecated
}

// valueIsStruct checks if the provided value is a struct or a struct pointer.
func valueIsStruct(value reflect.Value) bool {
	return value.Kind() == reflect.Struct ||
		(value.Kind() == reflect.Pointer && value.Type().Elem().Kind() == reflect.Struct)
}

// valueIsZero checks if a value is the zero value for its type
func valueIsZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return v.IsNil()
	default:
		return v.IsZero()
	}
}

// determineJsonNameForField determines the JSON field name for the given struct field.
func determineJsonNameForField(field reflect.StructField) string {
	// Get the JSON name for the field.
	jsonTag := field.Tag.Get("json")
	if jsonTag == "-" {
		// Field is explicitly excluded from JSON serialization.
		return ""
	}
	if jsonTag == "" {
		// No JSON tag, use the field name.
		return field.Name
	}

	// Parse json tag to get field name (before comma).
	jsonName := strings.Split(jsonTag, ",")[0]
	if jsonName == "" {
		// Empty name in tag (like `json:",omitempty"`) means use field name.
		return field.Name
	}

	return jsonName
}

func init() {
	templates.Register(check.Template{
		HumanName:   "Deprecated Field",
		Key:         templateKey,
		Description: "Checks for deprecated fields being set in custom resources",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.Any},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(_ params.Params) (check.Func, error) {
			return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				objValue := reflect.ValueOf(object.K8sObject)
				deprecatedFields := findDeprecatedFields(objValue, "")
				return deprecatedFieldsToDiagnostics(deprecatedFields)
			}, nil
		}),
	})
}

func deprecatedFieldsToDiagnostics(deprecatedFields []string) []diagnostic.Diagnostic {
	diagnostics := make([]diagnostic.Diagnostic, 0, len(deprecatedFields))
	for _, fieldPath := range deprecatedFields {
		diagnostics = append(diagnostics, diagnostic.Diagnostic{
			Message: fmt.Sprintf("field %q is deprecated and should not be set", fieldPath),
		})
	}
	return diagnostics
}
