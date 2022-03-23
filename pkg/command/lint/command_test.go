package lint

import (
	"github.com/spf13/cobra"
	"reflect"
	"testing"
)

func TestCommand_InvalidResources(t *testing.T) {
	tests := []struct {
		name       string
		cmd        *cobra.Command
		returnCode int
	}{
		{name: "InvalidResources", cmd: createLintCommand(), returnCode: 1},
		{name: "ValidResources", returnCode: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Command(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Command() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createLintCommand() *cobra.Command {

}
