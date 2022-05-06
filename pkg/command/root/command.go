package root

import (
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.stackrox.io/kube-linter/pkg/command/checks"
	"golang.stackrox.io/kube-linter/pkg/command/lint"
	"golang.stackrox.io/kube-linter/pkg/command/templates"
	"golang.stackrox.io/kube-linter/pkg/command/version"
)

const (
	colorFlag = "with-color"
)

// Command is the root command.
func Command() *cobra.Command {
	c := &cobra.Command{
		Use:           filepath.Base(os.Args[0]),
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			// Only forcefully set colorful output if the flag has been set.
			if cmd.Flags().Changed(colorFlag) {
				color.NoColor = false
			}
		},
	}
	c.AddCommand(
		checks.Command(),
		lint.Command(),
		templates.Command(),
		version.Command(),
	)
	c.PersistentFlags().Bool(colorFlag, true, "Force color output")
	return c
}
