package version

import (
	"fmt"

	"github.com/spf13/cobra"
	"golang.stackrox.io/kube-linter/internal/version"
)

// Command defines the version command
func Command() *cobra.Command {
	c := &cobra.Command{
		Use:   "version",
		Short: "Print version and exit",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Println(version.Get())
		},
	}
	return c
}
