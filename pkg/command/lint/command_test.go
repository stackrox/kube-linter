package lint

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/spf13/cobra"

	// Register templates
	_ "golang.stackrox.io/kube-linter/pkg/templates/all"
)

func TestCommand_InvalidResources(t *testing.T) {
	testDataPath := getTestDataDir()
	tests := []struct {
		name    string
		cmd     *cobra.Command
		failure bool
		output  string
	}{
		{name: "InvalidPodResource", cmd: createLintCommand(testDataPath+"invalid-pod-resources.yaml", "--fail-on-invalid-resource"), failure: true},
		{name: "InvalidPVCResource", cmd: createLintCommand(testDataPath+"invalid-pvc-resources.yaml", "--fail-on-invalid-resource"), failure: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.Execute()
			if err == nil && tt.failure {
				t.Fail()
			}
		})
	}
}

func createLintCommand(args ...string) *cobra.Command {
	c := Command()
	c.SetArgs(args)
	return c
}

func getTestDataDir() string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	return basepath + "/testdata/"
}
