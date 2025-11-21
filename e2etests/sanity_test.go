//go:build e2e

package e2etests

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	kubeLinterBinEnv = "KUBE_LINTER_BIN"
)

var (
	expectedOutRegex = regexp.MustCompile(`found \d+ lint errors`)
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

	gitCloneOut, err := exec.Command("git", "clone", "https://github.com/helm/charts.git", chartsDir).CombinedOutput()
	require.NoError(t, err, "Git clone failed. output: %s", string(gitCloneOut))

	kubeLinterOut, err := exec.Command(kubeLinterBin, "lint", chartsDir, "--config", "testdata/all-built-in-config.yaml").CombinedOutput()
	// Something will fail for sure, so kube-linter will not return a success code.
	require.Error(t, err)
	exitErr, ok := err.(*exec.ExitError)
	require.True(t, ok)
	outAsStr := string(kubeLinterOut)
	assert.Equal(t, 1, exitErr.ExitCode(), "unexpected exit code: %d; output from kube-linter: %v", exitErr.ExitCode(), outAsStr)
	assert.True(t, expectedOutRegex.MatchString(outAsStr), "unexpected output: %s", outAsStr)
}

func TestKubeLinterExitsWithNonZeroCodeOnEmptyDir(t *testing.T) {
	kubeLinterBin := os.Getenv(kubeLinterBinEnv)

	tmpDir, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(tmpDir))
	}()
	kubeLinterOut, err := exec.Command(kubeLinterBin, "lint", tmpDir, "--fail-if-no-objects-found").CombinedOutput()

	// error is expected
	require.Error(t, err)
	exitErr, ok := err.(*exec.ExitError)
	require.True(t, ok)
	outAsStr := string(kubeLinterOut)
	assert.Equal(t, 1, exitErr.ExitCode(), "unexpected exit code: %d; output from kube-linter: %v", exitErr.ExitCode(), outAsStr)
	msg := "no valid objects found"
	assert.True(t, strings.Contains(outAsStr, fmt.Sprintf("Error: %s", msg)), "unexpected output, it should contain: %s", outAsStr)

	// without the switch only warning is printed to stderr
	kubeLinterOut, err = exec.Command(kubeLinterBin, "lint", tmpDir).CombinedOutput()
	require.NoError(t, err)
	outAsStr = string(kubeLinterOut)
	assert.True(t, strings.Contains(outAsStr, fmt.Sprintf("Warning: %s", msg)), "unexpected output, it should contain: %s", outAsStr)
}
