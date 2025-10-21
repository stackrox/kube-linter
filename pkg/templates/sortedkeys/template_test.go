package sortedkeys

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
	appsV1 "k8s.io/api/apps/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/sortedkeys/internal/params"
)

func TestSortedKeys(t *testing.T) {
	suite.Run(t, new(SortedKeysTestSuite))
}

type SortedKeysTestSuite struct {
	templates.TemplateTestSuite
	ctx *mocks.MockLintContext
}

func (s *SortedKeysTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

type testCase struct {
	name                string
	yamlFile            string
	deploymentName      string
	recursive           bool
	expectedDiagnostics []string
}

// Helper function to load YAML from testdata
func (s *SortedKeysTestSuite) loadTestYAML(filename string) []byte {
	s.T().Helper()
	path := filepath.Join("testdata", filename)
	content, err := os.ReadFile(path) //nolint:gosec // Test file paths are controlled and safe
	s.Require().NoError(err, "Failed to read test file %s", filename)
	return content
}

func (s *SortedKeysTestSuite) TestSortedKeysTableDriven() {
	testCases := []testCase{
		{
			name:                "sorted_deployment",
			yamlFile:            "sorted-deployment.yaml",
			deploymentName:      "sorted-deployment",
			recursive:           true,
			expectedDiagnostics: nil,
		},
		{
			name:           "unsorted_top_level_keys",
			yamlFile:       "unsorted-top-level.yaml",
			deploymentName: "unsorted-deployment",
			recursive:      true,
			expectedDiagnostics: []string{
				"Keys are not sorted at root. Expected order: [apiVersion, kind, metadata, spec], got: [apiVersion, metadata, spec, kind]",
			},
		},
		{
			name:           "unsorted_nested_keys",
			yamlFile:       "unsorted-nested.yaml",
			deploymentName: "unsorted-nested",
			recursive:      true,
			expectedDiagnostics: []string{
				"Keys are not sorted at spec.template. Expected order: [metadata, spec], got: [spec, metadata]",
				"Keys are not sorted at spec.template.spec.containers[0]. Expected order: [image, name], got: [name, image]",
			},
		},
		{
			name:                "non_recursive_check",
			yamlFile:            "non-recursive.yaml",
			deploymentName:      "non-recursive",
			recursive:           false,
			expectedDiagnostics: nil,
		},
		{
			name:           "configmap_unsorted_data",
			yamlFile:       "configmap-unsorted-data.yaml",
			deploymentName: "unsorted-configmap",
			recursive:      true,
			expectedDiagnostics: []string{
				"Keys are not sorted at data. Expected order: [apple, zebra], got: [zebra, apple]",
			},
		},
		{
			name:           "multi_container",
			yamlFile:       "multi-container.yaml",
			deploymentName: "multi-container",
			recursive:      true,
			expectedDiagnostics: []string{
				"Keys are not sorted at spec.template.spec.containers[0]. Expected order: [image, name, ports], got: [name, image, ports]",
				"Keys are not sorted at spec.template.spec.containers[0].ports[0]. Expected order: [containerPort, name, protocol], got: [protocol, containerPort, name]",
				"Keys are not sorted at spec.template.spec.containers[1]. Expected order: [image, name, volumeMounts], got: [name, image, volumeMounts]",
				"Keys are not sorted at spec.template.spec.containers[1].volumeMounts[0]. Expected order: [mountPath, name], got: [name, mountPath]",
				"Keys are not sorted at spec.template.spec.volumes[0]. Expected order: [emptyDir, name], got: [name, emptyDir]",
			},
		},
		{
			name:                "edge_cases",
			yamlFile:            "edge-cases.yaml",
			deploymentName:      "edge-cases",
			recursive:           true,
			expectedDiagnostics: nil,
		},
		{
			name:           "numeric_keys",
			yamlFile:       "numeric-keys.yaml",
			deploymentName: "numeric-keys",
			recursive:      true,
			expectedDiagnostics: []string{
				"Keys are not sorted at data. Expected order: [1-config, 10-config, 2-config, config-a, config-z], got: [10-config, 2-config, 1-config, config-z, config-a]",
			},
		},
		{
			name:           "complex_service",
			yamlFile:       "complex-service.yaml",
			deploymentName: "complex-service",
			recursive:      true,
			expectedDiagnostics: []string{
				"Keys are not sorted at spec.ports[0]. Expected order: [name, port, protocol, targetPort], got: [port, targetPort, protocol, name]",
			},
		},
		{
			name:           "mixed_sorting",
			yamlFile:       "mixed-sorting.yaml",
			deploymentName: "mixed-sorting",
			recursive:      true,
			expectedDiagnostics: []string{
				"Keys are not sorted at spec.selector.matchLabels. Expected order: [app, version], got: [version, app]",
			},
		},

		// Complex test cases from testdata files
		{
			name:           "deeply_nested_unsorted",
			yamlFile:       "deeply-nested-unsorted.yaml",
			deploymentName: "deeply-nested",
			recursive:      true,
			expectedDiagnostics: []string{
				"Keys are not sorted at metadata.annotations. Expected order: [app.kubernetes.io/version, kubernetes.io/change-cause], got: [kubernetes.io/change-cause, app.kubernetes.io/version]",
				"Keys are not sorted at spec.template.spec.containers[0].resources. Expected order: [limits, requests], got: [requests, limits]",
				"Keys are not sorted at spec.template.spec.containers[0].resources.requests. Expected order: [cpu, memory], got: [memory, cpu]",
				"Keys are not sorted at spec.template.spec.containers[0].resources.limits. Expected order: [cpu, memory], got: [memory, cpu]",
			},
		},

		{
			name:           "complex_pod_spec",
			yamlFile:       "complex-pod-spec.yaml",
			deploymentName: "complex-pod-spec",
			recursive:      true,
			expectedDiagnostics: []string{
				"Keys are not sorted at spec.template.spec. Expected order: [containers, initContainers, securityContext, serviceAccountName, volumes], got: [containers, initContainers, serviceAccountName, securityContext, volumes]",
				"Keys are not sorted at spec.template.spec.containers[0]. Expected order: [image, name, securityContext, volumeMounts], got: [image, name, volumeMounts, securityContext]",
				"Keys are not sorted at spec.template.spec.containers[0].securityContext. Expected order: [allowPrivilegeEscalation, runAsNonRoot], got: [runAsNonRoot, allowPrivilegeEscalation]",
				"Keys are not sorted at spec.template.spec.securityContext. Expected order: [fsGroup, runAsGroup, runAsUser], got: [fsGroup, runAsUser, runAsGroup]",
			},
		},

		{
			name:                "all_sorted_complex",
			yamlFile:            "all-sorted-complex.yaml",
			deploymentName:      "all-sorted-complex",
			recursive:           true,
			expectedDiagnostics: nil,
		},

		// YAML reference and merge key test cases
		{
			name:                "reused_labels_sorted",
			yamlFile:            "reused-labels-sorted.yaml",
			deploymentName:      "anchor-deployment",
			recursive:           true,
			expectedDiagnostics: nil, // All keys are sorted
		},

		{
			name:           "reused_labels_unsorted",
			yamlFile:       "reused-labels-unsorted.yaml",
			deploymentName: "anchor-unsorted",
			recursive:      true,
			expectedDiagnostics: []string{
				"Keys are not sorted at metadata.labels. Expected order: [app, environment, version], got: [version, app, environment]",
				"Keys are not sorted at spec.template.spec.containers[0]. Expected order: [image, name], got: [name, image]",
			},
		},

		{
			name:           "container_with_merge",
			yamlFile:       "container-with-merge.yaml",
			deploymentName: "complex-anchors",
			recursive:      true,
			expectedDiagnostics: []string{
				"Keys are not sorted at spec.template.spec.containers[0]. Expected order: [env, image, name, resources], got: [image, name, resources, env]",
				"Keys are not sorted at spec.template.spec.containers[1]. Expected order: [<<, env, image, name, resources], got: [<<, name, image, resources, env]",
				"Keys are not sorted at spec.template.spec.containers[2]. Expected order: [<<, env, image, name], got: [<<, name, env, image]",
			},
		},

		{
			name:           "configmap_merge_unsorted",
			yamlFile:       "configmap-merge-unsorted.yaml",
			deploymentName: "config-with-anchors",
			recursive:      true,
			expectedDiagnostics: []string{
				"Keys are not sorted at data.base-config. Expected order: [retries, timeout], got: [timeout, retries]",
				"Keys are not sorted at data.service-a-config. Expected order: [<<, auth, endpoint], got: [<<, endpoint, auth]",
				"Keys are not sorted at data.service-b-config. Expected order: [<<, aa-priority, zz-custom], got: [<<, zz-custom, aa-priority]",
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Reset context for each test case
			s.ctx = mocks.NewMockContext()

			// Load YAML content from testdata
			yamlContent := s.loadTestYAML(tc.yamlFile)

			// Create deployment object (using generic deployment for all test cases)
			deployment := &appsV1.Deployment{
				TypeMeta: metaV1.TypeMeta{
					APIVersion: "apps/v1",
					Kind:       "Deployment",
				},
				ObjectMeta: metaV1.ObjectMeta{
					Name: tc.deploymentName,
				},
			}

			// Add object with raw YAML
			s.ctx.AddObjectWithRaw(tc.deploymentName, deployment, yamlContent)

			// Build expected diagnostics map
			expectedDiags := map[string][]diagnostic.Diagnostic{
				tc.deploymentName: nil,
			}
			if len(tc.expectedDiagnostics) > 0 {
				diags := make([]diagnostic.Diagnostic, 0, len(tc.expectedDiagnostics))
				for _, msg := range tc.expectedDiagnostics {
					diags = append(diags, diagnostic.Diagnostic{Message: msg})
				}
				expectedDiags[tc.deploymentName] = diags
			}

			// Validate
			s.Validate(s.ctx, []templates.TestCase{
				{
					Param: params.Params{
						Recursive: tc.recursive,
					},
					Diagnostics:              expectedDiags,
					ExpectInstantiationError: false,
				},
			})
		})
	}
}
