package sortedkeys

import (
	"fmt"
	"sort"
	"strings"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/sortedkeys/internal/params"
	"gopkg.in/yaml.v3"
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
			return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				// Parse the raw YAML to preserve key order
				var node yaml.Node
				if err := yaml.Unmarshal(object.Metadata.Raw, &node); err != nil {
					// Skip objects that can't be parsed
					return nil
				}

				var diagnostics []diagnostic.Diagnostic

				// The root node is a document node, we need to check its content
				if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
					diagnostics = checkNode(node.Content[0], "", p.Recursive)
				} else if node.Kind == yaml.MappingNode {
					diagnostics = checkNode(&node, "", p.Recursive)
				}

				return diagnostics
			}, nil
		}),
	})
}

// checkNode recursively checks if keys in a YAML node are sorted
func checkNode(node *yaml.Node, path string, recursive bool) []diagnostic.Diagnostic {
	if node == nil {
		return nil
	}

	var diagnostics []diagnostic.Diagnostic

	switch node.Kind {
	case yaml.MappingNode:
		// Extract keys from the mapping node
		// In yaml.v3, mapping nodes store key-value pairs as alternating elements in Content
		var keys []string
		keyPositions := make(map[string]int) // Track original positions

		for i := 0; i < len(node.Content); i += 2 {
			if node.Content[i].Kind == yaml.ScalarNode {
				key := node.Content[i].Value
				keys = append(keys, key)
				keyPositions[key] = i / 2
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
		if recursive {
			for i := 0; i < len(node.Content); i += 2 {
				keyNode := node.Content[i]
				valueNode := node.Content[i+1]

				childPath := path
				if childPath == "" {
					childPath = keyNode.Value
				} else {
					childPath = path + "." + keyNode.Value
				}

				childDiagnostics := checkNode(valueNode, childPath, recursive)
				diagnostics = append(diagnostics, childDiagnostics...)
			}
		}

	case yaml.SequenceNode:
		// For sequences, check each element if it's a mapping
		if recursive {
			for idx, item := range node.Content {
				childPath := fmt.Sprintf("%s[%d]", path, idx)
				childDiagnostics := checkNode(item, childPath, recursive)
				diagnostics = append(diagnostics, childDiagnostics...)
			}
		}
	}

	return diagnostics
}
