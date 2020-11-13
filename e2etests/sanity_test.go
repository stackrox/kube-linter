// +build e2e

package e2etests

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	kubeLinterBinEnv = "KUBE_LINTER_BIN"
)

func TestKubeLinterWithBuiltInChecksDoesntCrashOnHelmChartsRepo(t *testing.T) {
	kubeLinterBin := os.Getenv(kubeLinterBinEnv)
	require.NotEmpty(t, kubeLinterBin, "Please set %s", kubeLinterBinEnv)

	_, err := os.Stat(kubeLinterBin)
	require.NoError(t, err)

	tmpDir, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(tmpDir))
	}()

	chartsDir := filepath.Join(tmpDir, "charts")

	gitCloneOut, err := exec.Command("git", "clone", "git@github.com:helm/charts.git", chartsDir).CombinedOutput()
	require.NoError(t, err, "Git clone failed. output: %s", string(gitCloneOut))

	kubeLinterOut, err := exec.Command(kubeLinterBin, "lint", chartsDir, "--config", "testdata/all-built-in-config.yaml").CombinedOutput()
	// Something will fail for sure, so kube-linter will not return a success code.
	require.Error(t, err)
	exitErr, ok := err.(*exec.ExitError)
	require.True(t, ok)
	assert.Equal(t, 1, exitErr.ExitCode(), "unexpected exit code: %d; output from kube-linter: %v", exitErr.ExitCode(), string(kubeLinterOut))
}
