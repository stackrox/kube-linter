//go:build e2e

package e2etests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestKubeChainsawIntegration verifies that the kube-chainsaw template
// integrates correctly with kube-linter and detects RBAC issues.
func TestKubeChainsawIntegration(t *testing.T) {
	kubeLinterBin := os.Getenv(kubeLinterBinEnv)
	require.NotEmpty(t, kubeLinterBin, "Please set %s", kubeLinterBinEnv)

	_, err := os.Stat(kubeLinterBin)
	require.NoError(t, err)

	// Test fixture with ServiceAccount, ClusterRoleBinding to cluster-admin, and Pod
	fixtureDir := filepath.Join("..", "tests", "checks", "kube-chainsaw")

	// Run kube-linter with kube-chainsaw-rbac check enabled
	kubeLinterOut, err := exec.Command(
		kubeLinterBin,
		"lint",
		fixtureDir,
		"--include", "kube-chainsaw-rbac",
	).CombinedOutput()

	// kube-linter should exit with error code 1 when findings are detected
	require.Error(t, err)
	exitErr, ok := err.(*exec.ExitError)
	require.True(t, ok)
	outAsStr := string(kubeLinterOut)

	assert.Equal(t, 1, exitErr.ExitCode(),
		"unexpected exit code: %d; output from kube-linter: %v",
		exitErr.ExitCode(), outAsStr)

	// Verify that cluster-admin finding is detected
	assert.Contains(t, outAsStr, "cluster-admin privileges",
		"expected cluster-admin privileges finding in output: %s", outAsStr)

	// Verify severity is critical
	assert.True(t, strings.Contains(outAsStr, "CRITICAL") || strings.Contains(outAsStr, "critical"),
		"expected 'critical' severity in output: %s", outAsStr)

	// Verify the finding references the dangerous pod
	assert.Contains(t, outAsStr, "privileged-pod",
		"expected finding to reference 'privileged-pod': %s", outAsStr)

	// Verify the kube-chainsaw-rbac check ran
	assert.Contains(t, outAsStr, "kube-chainsaw-rbac",
		"expected kube-chainsaw-rbac check in output: %s", outAsStr)
}
