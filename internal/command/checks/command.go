package checks

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.stackrox.io/kube-linter/internal/builtinchecks"
	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/command/common"
	"golang.stackrox.io/kube-linter/internal/defaultchecks"
	"golang.stackrox.io/kube-linter/internal/ternary"
)

var (
	dashes = func() string {
		var sb strings.Builder
		for i := 0; i < 30; i++ {
			sb.WriteRune('-')
		}
		return sb.String()
	}()

	formatsToRenderFuncs = map[string]func([]check.Check, io.Writer){
		common.PlainFormat: func(checks []check.Check, out io.Writer) {
			for i, check := range checks {
				fmt.Fprintf(out, "Name: %s\nDescription: %s\nTemplate: %s\nParameters: %v\nEnabled by default: %v\n",
					check.Name, check.Description, check.Template, check.Params, defaultchecks.List.Contains(check.Name))
				if i != len(checks)-1 {
					fmt.Fprintf(out, "\n%s\n\n", dashes)
				}
			}
		},
		common.MarkdownFormat: func(checks []check.Check, out io.Writer) {
			fmt.Fprintf(out, "The following table enumerates built-in checks:\n")
			fmt.Fprintf(out, "\n\n| Name | Enabled by default | Description | Template | Parameters |\n --- | --- | --- | --- | --- | \n")
			for _, check := range checks {
				var params string
				if len(check.Params) == 0 {
					params = "none"
				} else {
					var sb strings.Builder
					for key, value := range check.Params {
						sb.WriteString(fmt.Sprintf("- `%s`: `%s` <br />", key, value))
					}
					params = sb.String()
				}
				fmt.Fprintf(out, "|`%s`|%s|%s|%s|%s|\n",
					check.Name,
					ternary.String(defaultchecks.List.Contains(check.Name), "Yes", "No"),
					check.Description, check.Template, params)
			}
		},
	}
)

func listCommand() *cobra.Command {
	format := common.FormatWrapper{Format: common.PlainFormat}
	c := &cobra.Command{
		Use:   "list",
		Short: "list built-in checks",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			checks, err := builtinchecks.List()
			if err != nil {
				return err
			}
			renderFunc := formatsToRenderFuncs[format.Format]
			if renderFunc == nil {
				return errors.Errorf("unknown format: %q", format.Format)
			}
			renderFunc(checks, os.Stdout)
			return nil
		},
	}
	c.Flags().Var(&format, "format", "output format")
	return c
}

// Command defines the root of the checks command.
func Command() *cobra.Command {
	c := &cobra.Command{
		Use: "checks",
	}
	c.AddCommand(listCommand())
	return c
}
