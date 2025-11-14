package lint

import (
	"testing"

	"github.com/spf13/cobra"

	// Register templates
	_ "golang.stackrox.io/kube-linter/pkg/templates/all"
)

func TestCommand_InvalidResources(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *cobra.Command
		failure bool
		output  string
	}{
		{name: "InvalidYAML", cmd: createLintCommand("./testdata/invalid.yaml", "--fail-on-invalid-resource"), failure: true},
		{name: "NonexistentFile", cmd: createLintCommand("./testdata/foo-bar.yaml", "--fail-on-invalid-resource"), failure: true},
		{name: "ValidPod", cmd: createLintCommand("./testdata/valid-pod.yaml", "--fail-on-invalid-resource"), failure: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.Execute()
			if err == nil && tt.failure || err != nil && !tt.failure {
				t.Fail()
			}
		})
	}
}

func createLintCommand(args ...string) *cobra.Command {
	c := Command()
	c.SilenceUsage = true
	c.SetArgs(args)
	return c
}
