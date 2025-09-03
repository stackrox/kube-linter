package cel

import (
	"encoding/json"
	"fmt"

	"github.com/google/cel-go/cel"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/cel/internal/params"
)

const (
	templateKey = "cel-expression"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "CEL",
		Key:         templateKey,
		Description: "Flag objects with CEL expression",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.Any},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			return func(ctx lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				msg, err := evaluate(p.Check, object, ctx.Objects())
				if err != nil {
					return []diagnostic.Diagnostic{
						{Message: fmt.Sprintf("error evaluating CEL check expression: %v", err)},
					}
				}
				if msg != "" {
					return []diagnostic.Diagnostic{
						{Message: fmt.Sprintf("CEL check expression returned: %v", msg)},
					}
				}
				return nil
			}, nil
		}),
	})
}

func evaluate(check string, object lintcontext.Object, objects []lintcontext.Object) (string, error) {
	// Convert object to map via JSON marshaling/unmarshaling for CEL compatibility
	// We need to marshal the underlying K8sObject, not the lintcontext.Object
	objectMap, err := toMap(object.K8sObject)
	if err != nil {
		return "", fmt.Errorf("failed to convert object to map: %w", err)
	}

	// Convert objects to maps via JSON marshaling/unmarshaling
	objectsMaps := make([]map[string]any, len(objects))
	for i, obj := range objects {
		objMap, err := toMap(obj.K8sObject)
		if err != nil {
			return "", fmt.Errorf("failed to convert object %s to map: %w", obj.GetK8sObjectName().String(), err)
		}
		objectsMaps[i] = objMap
	}

	e, err := cel.NewEnv(
		cel.Variable("object", cel.MapType(cel.StringType, cel.AnyType)),
		cel.Variable("objects", cel.ListType(cel.MapType(cel.StringType, cel.AnyType))),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create CEL environment: %w", err)
	}
	ast, iss := e.Compile(check)
	if iss.Err() != nil {
		return "", fmt.Errorf("failed to compile CEL expression: %w", iss.Err())
	}
	prg, err := e.Program(ast)
	if err != nil {
		return "", fmt.Errorf("failed to create CEL program: %w", err)
	}
	out, _, err := prg.Eval(map[string]any{
		"object":  objectMap,
		"objects": objectsMaps,
	})
	if err != nil {
		return "", fmt.Errorf("failed to evaluate CEL expression: %w", err)
	}

	o, ok := out.Value().(string)
	if !ok {
		return "", fmt.Errorf("expected string, got %v", out.Value())
	}
	return o, nil
}

func toMap(obj any) (map[string]any, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal object: %w", err)
	}
	var output map[string]any
	if err := json.Unmarshal(bytes, &output); err != nil {
		return nil, fmt.Errorf("failed to unmarshal object: %w", err)
	}
	return output, nil
}
