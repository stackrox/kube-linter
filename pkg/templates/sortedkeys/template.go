package sortedkeys

import (
	"fmt"
	"sort"
	"strings"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/sortedkeys/internal/params"
)

const templateKey = "sorted-keys"

func init() {
	templates.Register(check.Template{
		HumanName:   "Sorted Keys",
		Key:         templateKey,
		Description: "Flag YAML keys that are not sorted in alphabetical order",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.Any},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			ignoredKeys := make(map[string]struct{}, len(p.IgnoredKeys))
			for _, key := range p.IgnoredKeys {
				ignoredKeys[key] = struct{}{}
			}
			return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				// Parse the raw YAML to preserve key order
				file, err := parser.ParseBytes(object.Metadata.Raw, parser.ParseComments)
				if err != nil {
					// Skip objects that can't be parsed
					return nil
				}

				var diagnostics []diagnostic.Diagnostic
				c := checker{recursive: p.Recursive, ignoredKeys: ignoredKeys}

				// Check all documents in the file
				for _, doc := range file.Docs {
					if doc != nil && doc.Body != nil {
						diagnostics = append(diagnostics, c.checkNode(doc.Body, "")...)
					}
				}

				return diagnostics
			}, nil
		}),
	})
}

type checker struct {
	recursive   bool
	ignoredKeys map[string]struct{}
}

// checkNode recursively checks if keys in a YAML node are sorted
func (c checker) checkNode(node ast.Node, path string) []diagnostic.Diagnostic {
	if node == nil {
		return nil
	}

	var diagnostics []diagnostic.Diagnostic

	switch n := node.(type) {
	case *ast.MappingNode:
		// MappingNode contains Values which are MappingValueNodes
		var keys []string
		keyPositions := make(map[string]int)

		for i, value := range n.Values {
			// Values in a MappingNode are already MappingValueNode pointers
			// Extract the key
			key := getKeyString(value.Key)
			if key != "" {
				keys = append(keys, key)
				keyPositions[key] = i
			}
		}

		// Check if keys are sorted
		if len(keys) > 1 {
			sortedKeys := make([]string, len(keys))
			copy(sortedKeys, keys)
			sort.Strings(sortedKeys)

			// Find the first key that is out of order
			for i := 0; i < len(keys); i++ {
				if keys[i] != sortedKeys[i] {
					location := path
					if location == "" {
						location = "root"
					}

					diagnostics = append(diagnostics, diagnostic.Diagnostic{
						Message: fmt.Sprintf(
							"Keys are not sorted at %s. Expected order: [%s], got: [%s]",
							location,
							strings.Join(sortedKeys, ", "),
							strings.Join(keys, ", "),
						),
					})
					// Only report once per level
					break
				}
			}
		}

		// Recursively check child nodes if recursive is enabled
		if c.recursive {
			for _, value := range n.Values {
				key := getKeyString(value.Key)
				if _, ignored := c.ignoredKeys[key]; ignored {
					continue
				}
				childPath := path
				if childPath == "" {
					childPath = key
				} else {
					childPath = path + "." + key
				}

				childDiagnostics := c.checkNode(value.Value, childPath)
				diagnostics = append(diagnostics, childDiagnostics...)
			}
		}

	case *ast.SequenceNode:
		// For sequences, check each element if it's a mapping
		if c.recursive {
			for idx, item := range n.Values {
				childPath := fmt.Sprintf("%s[%d]", path, idx)
				childDiagnostics := c.checkNode(item, childPath)
				diagnostics = append(diagnostics, childDiagnostics...)
			}
		}

	case *ast.AnchorNode:
		// Handle anchor nodes by checking their value
		diagnostics = append(diagnostics, c.checkNode(n.Value, path)...)

	case *ast.AliasNode:
		// Skip alias nodes - they reference already checked content
		// No need to check as the original anchor was already checked

	case *ast.MergeKeyNode:
		// Handle merge keys (<<: *alias)
		// The merge key itself is represented as a special key
		// No special handling needed as it will be treated as a key
	}

	return diagnostics
}

// getKeyString extracts the string representation of a key node
func getKeyString(node ast.Node) string {
	switch n := node.(type) {
	case *ast.StringNode:
		return n.Value
	case *ast.IntegerNode:
		return n.String()
	case *ast.FloatNode:
		return n.String()
	case *ast.BoolNode:
		if n.Value {
			return "true"
		}
		return "false"
	case *ast.MergeKeyNode:
		return "<<"
	case *ast.NullNode:
		return "null"
	default:
		return ""
	}
}
