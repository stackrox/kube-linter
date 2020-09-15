package root

import (
	"os"

	"github.com/spf13/cobra"
	"golang.stackrox.io/kube-linter/internal/command/lint"
)

// Command is the root command.
func Command() *cobra.Command {
	c := &cobra.Command{
		Use:          os.Args[0],
		SilenceUsage: true,
	}
	c.AddCommand(lint.Command())
	return c
}
