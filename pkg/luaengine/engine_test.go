package luaengine

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.stackrox.io/kube-linter/pkg/k8sutil"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestEngine_ExecuteCheck(t *testing.T) {
	tests := []struct {
		name          string
		script        string
		object        runtime.Object
		expectDiags   int
		expectError   bool
		expectMessage string
	}{
		{
			name: "simple check returning no diagnostics",
			script: `
				function check()
					return {}
				end
			`,
			object: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-pod",
				},
			},
			expectDiags: 0,
		},
		{
			name: "check returning one diagnostic",
			script: `
				function check()
					return {
						diagnostic("Test diagnostic message")
					}
				end
			`,
			object: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-pod",
				},
			},
			expectDiags:   1,
			expectMessage: "Test diagnostic message",
		},
		{
			name: "check using extract.labels",
			script: `
				function check()
					local diagnostics = {}
					local labels = extract.labels()
					if labels and labels.environment == "production" then
						table.insert(diagnostics, diagnostic("Found production label"))
					end
					return diagnostics
				end
			`,
			object: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-pod",
					Labels: map[string]string{
						"environment": "production",
					},
				},
			},
			expectDiags:   1,
			expectMessage: "Found production label",
		},
		{
			name: "check using extract.podSpec",
			script: `
				function check()
					local diagnostics = {}
					local podSpec = extract.podSpec()
					if podSpec and podSpec.containers then
						for i, container in ipairs(podSpec.containers) do
							if container.image and string.match(container.image, ":latest$") then
								table.insert(diagnostics, diagnostic("Found latest tag"))
							end
						end
					end
					return diagnostics
				end
			`,
			object: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-pod",
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "test-container",
							Image: "nginx:latest",
						},
					},
				},
			},
			expectDiags:   1,
			expectMessage: "Found latest tag",
		},
		{
			name: "script with syntax error",
			script: `
				function check(
					-- missing closing parenthesis
				end
			`,
			object: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-pod",
				},
			},
			expectError: true,
		},
		{
			name: "script with runtime error",
			script: `
				function check()
					error("Runtime error")
				end
			`,
			object: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-pod",
				},
			},
			expectError: true,
		},
		{
			name: "check with nil return",
			script: `
				function check()
					return nil
				end
			`,
			object: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-pod",
				},
			},
			expectDiags: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := New(tt.script, 5*time.Second)

			// Create lint context and object
			lintCtx := &mockLintContext{} // Empty context for tests
			obj := lintcontext.Object{
				K8sObject: tt.object.(k8sutil.Object),
			}

			diagnostics, err := engine.ExecuteCheck(lintCtx, obj)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, diagnostics, tt.expectDiags)

			if tt.expectDiags > 0 && tt.expectMessage != "" {
				assert.Equal(t, tt.expectMessage, diagnostics[0].Message)
			}
		})
	}
}

func TestEngine_Timeout(t *testing.T) {
	script := `
		function check()
			-- Infinite loop to test timeout
			while true do
				-- busy wait
			end
			return {}
		end
	`

	engine := New(script, 100*time.Millisecond)

	lintCtx := &mockLintContext{}
	obj := lintcontext.Object{
		K8sObject: &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-pod",
			},
		},
	}

	start := time.Now()
	_, err := engine.ExecuteCheck(lintCtx, obj)
	elapsed := time.Since(start)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "timeout")
	assert.Less(t, elapsed, 500*time.Millisecond, "Should timeout quickly")
}

func TestEngine_DefaultTimeout(t *testing.T) {
	engine := New("function check() return {} end", 0)
	assert.Equal(t, defaultTimeout, engine.timeout)
}

func TestEngine_MaxTimeout(t *testing.T) {
	engine := New("function check() return {} end", 60*time.Second)
	assert.Equal(t, maxTimeout, engine.timeout)
}

// Mock implementations for testing
type mockLintContext struct{}

func (m *mockLintContext) Objects() []lintcontext.Object               { return nil }
func (m *mockLintContext) InvalidObjects() []lintcontext.InvalidObject { return nil }
