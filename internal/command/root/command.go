package root

import (
	"os"

	"github.com/spf13/cobra"
	"golang.stackrox.io/kube-linter/internal/command/checks"
	"golang.stackrox.io/kube-linter/internal/command/lint"
	"golang.stackrox.io/kube-linter/internal/command/templates"
)

// Command is the root command.
func Command() *cobra.Command {
	c := &cobra.Command{
		Use:           os.Args[0],
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	c.AddCommand(
		checks.Command(),
		lint.Command(),
		templates.Command(),
	)
	return c
}
