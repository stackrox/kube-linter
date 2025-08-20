package luascript

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/instantiatedcheck"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/templates"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestLuaScriptTemplate(t *testing.T) {
	// Test that the template is registered
	template, found := templates.Get("lua-script")
	require.True(t, found, "lua-script template should be registered")
	assert.Equal(t, "Lua Script", template.HumanName)
	assert.Equal(t, "lua-script", template.Key)
}

func TestLuaScriptTemplate_InlineScript(t *testing.T) {
	check := &config.Check{
		Name:        "test-lua-inline",
		Description: "Test inline Lua script",
		Template:    "lua-script",
		Params: map[string]interface{}{
			"inline": `
				function check()
					local diagnostics = {}
					local labels = extract.labels()
					if labels and labels.env == "test" then
						table.insert(diagnostics, diagnostic("Found test environment"))
					end
					return diagnostics
				end
			`,
			"timeout": 5,
		},
	}

	instantiated, err := instantiatedcheck.ValidateAndInstantiate(check)
	require.NoError(t, err)

	// Test with matching labels
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-pod",
			Labels: map[string]string{
				"env": "test",
			},
		},
	}

	lintCtx := &mockLintContext{}
	obj := lintcontext.Object{K8sObject: pod}

	diagnostics := instantiated.Func(lintCtx, obj)
	require.Len(t, diagnostics, 1)
	assert.Equal(t, "Found test environment", diagnostics[0].Message)

	// Test with non-matching labels
	pod.Labels["env"] = "production"
	diagnostics = instantiated.Func(lintCtx, obj)
	assert.Len(t, diagnostics, 0)
}

func TestLuaScriptTemplate_FileScript(t *testing.T) {
	// Create temporary script file
	tempDir := t.TempDir()
	scriptPath := filepath.Join(tempDir, "test-script.lua")

	scriptContent := `
		function check()
			local diagnostics = {}
			local podSpec = extract.podSpec()
			if podSpec and podSpec.containers then
				for i, container in ipairs(podSpec.containers) do
					if container.image and string.match(container.image, ":latest$") then
						table.insert(diagnostics, diagnostic("Container uses latest tag"))
					end
				end
			end
			return diagnostics
		end
	`

	err := os.WriteFile(scriptPath, []byte(scriptContent), 0600)
	require.NoError(t, err)

	check := &config.Check{
		Name:        "test-lua-file",
		Description: "Test file-based Lua script",
		Template:    "lua-script",
		Params: map[string]interface{}{
			"script":  scriptPath,
			"timeout": 5,
		},
	}

	instantiated, err := instantiatedcheck.ValidateAndInstantiate(check)
	require.NoError(t, err)

	// Test with pod containing latest tag
	pod := &v1.Pod{
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
	}

	lintCtx := &mockLintContext{}
	obj := lintcontext.Object{K8sObject: pod}

	diagnostics := instantiated.Func(lintCtx, obj)
	require.Len(t, diagnostics, 1)
	assert.Equal(t, "Container uses latest tag", diagnostics[0].Message)
}

func TestLuaScriptTemplate_ValidationErrors(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]interface{}
		expectError string
	}{
		{
			name:        "missing script and inline",
			params:      map[string]interface{}{},
			expectError: "must specify either 'script' or 'inline' parameter",
		},
		{
			name: "both script and inline",
			params: map[string]interface{}{
				"script": "test.lua",
				"inline": "function check() end",
			},
			expectError: "cannot specify both 'script' and 'inline' parameters",
		},
		{
			name: "nonexistent script file",
			params: map[string]interface{}{
				"script": "/nonexistent/path/script.lua",
			},
			expectError: "reading script file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			check := &config.Check{
				Name:        "test-validation",
				Description: "Test validation",
				Template:    "lua-script",
				Params:      tt.params,
			}

			_, err := instantiatedcheck.ValidateAndInstantiate(check)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectError)
		})
	}
}

func TestLuaScriptTemplate_ScriptError(t *testing.T) {
	check := &config.Check{
		Name:        "test-script-error",
		Description: "Test script with error",
		Template:    "lua-script",
		Params: map[string]interface{}{
			"inline": `
				function check()
					error("Test error from Lua")
				end
			`,
		},
	}

	instantiated, err := instantiatedcheck.ValidateAndInstantiate(check)
	require.NoError(t, err)

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-pod",
		},
	}

	lintCtx := &mockLintContext{}
	obj := lintcontext.Object{K8sObject: pod}

	diagnostics := instantiated.Func(lintCtx, obj)
	require.Len(t, diagnostics, 1)
	assert.Contains(t, diagnostics[0].Message, "Lua script error")
}

func TestLuaScriptTemplate_ComplexCheck(t *testing.T) {
	check := &config.Check{
		Name:        "test-complex-check",
		Description: "Test complex Lua check",
		Template:    "lua-script",
		Params: map[string]interface{}{
			"inline": `
				function check()
					local diagnostics = {}
					local podSpec = extract.podSpec()
					
					if not podSpec then
						return diagnostics
					end
					
					-- Check for privileged containers
					local hasPrivileged = false
					if podSpec.containers then
						for i, container in ipairs(podSpec.containers) do
							if container.securityContext and 
							   container.securityContext.privileged == true then
								hasPrivileged = true
								break
							end
						end
					end
					
					-- Check for host networking
					local usesHostNetwork = podSpec.hostNetwork == true
					
					if hasPrivileged and usesHostNetwork then
						table.insert(diagnostics, diagnostic(
							"Pod uses both privileged containers and host networking"
						))
					end
					
					return diagnostics
				end
			`,
		},
	}

	instantiated, err := instantiatedcheck.ValidateAndInstantiate(check)
	require.NoError(t, err)

	// Test case 1: Pod with both privileged container and host networking
	pod1 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "privileged-hostnet-pod",
		},
		Spec: v1.PodSpec{
			HostNetwork: true,
			Containers: []v1.Container{
				{
					Name:  "privileged-container",
					Image: "nginx",
					SecurityContext: &v1.SecurityContext{
						Privileged: &[]bool{true}[0],
					},
				},
			},
		},
	}

	lintCtx := &mockLintContext{}
	obj := lintcontext.Object{K8sObject: pod1}

	diagnostics := instantiated.Func(lintCtx, obj)
	require.Len(t, diagnostics, 1)
	assert.Contains(t, diagnostics[0].Message, "privileged containers and host networking")

	// Test case 2: Pod with only privileged container (should pass)
	pod2 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "privileged-only-pod",
		},
		Spec: v1.PodSpec{
			HostNetwork: false,
			Containers: []v1.Container{
				{
					Name:  "privileged-container",
					Image: "nginx",
					SecurityContext: &v1.SecurityContext{
						Privileged: &[]bool{true}[0],
					},
				},
			},
		},
	}

	obj = lintcontext.Object{K8sObject: pod2}
	diagnostics = instantiated.Func(lintCtx, obj)
	assert.Len(t, diagnostics, 0)

	// Test case 3: Pod with only host networking (should pass)
	pod3 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "hostnet-only-pod",
		},
		Spec: v1.PodSpec{
			HostNetwork: true,
			Containers: []v1.Container{
				{
					Name:  "normal-container",
					Image: "nginx",
				},
			},
		},
	}

	obj = lintcontext.Object{K8sObject: pod3}
	diagnostics = instantiated.Func(lintCtx, obj)
	assert.Len(t, diagnostics, 0)
}

// Mock implementations for testing
type mockLintContext struct{}

func (m *mockLintContext) Objects() []lintcontext.Object               { return nil }
func (m *mockLintContext) InvalidObjects() []lintcontext.InvalidObject { return nil }
