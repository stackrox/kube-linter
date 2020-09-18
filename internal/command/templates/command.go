package templates

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/command/common"
	"golang.stackrox.io/kube-linter/internal/templates"
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

	formatsToRenderFuncs = map[string]func([]check.Template, io.Writer){
		common.PlainFormat: func(templates []check.Template, out io.Writer) {
			for i, template := range templates {
				fmt.Fprintf(out, "Name: %s\nDescription: %s\nSupported Objects: %v\n", template.Name, template.Description, template.SupportedObjectKinds.ObjectKinds)
				if len(template.Parameters) == 0 {
					fmt.Fprintln(out, "Parameters: none")
				} else {
					fmt.Fprintf(out, "Parameters:\n")
					for _, param := range template.Parameters {
						fmt.Fprintf(out, "\t%s:\n\t\tDescription: %s\n\t\tRequired: %v\n", param.ParamName, param.Description, param.Required)
					}
				}
				if i != len(templates)-1 {
					fmt.Fprintf(out, "\n%s\n\n", dashes)
				}
			}
		},
		common.MarkdownFormat: func(templates []check.Template, out io.Writer) {
			fmt.Fprintf(out, "The following table enumerates supported check templates:\n")
			fmt.Fprintf(out, "\n\n| Name | Description | Supported Objects | Parameters |\n --- | --- | --- | --- |\n")
			for _, template := range templates {
				var params string
				if len(template.Parameters) == 0 {
					params = "none"
				} else {
					var sb strings.Builder
					for _, param := range template.Parameters {
						sb.WriteString(fmt.Sprintf("- `%s`%s: %s <br />", param.ParamName, ternary.String(param.Required, " (required)", ""), param.Description))
					}
					params = sb.String()
				}
				fmt.Fprintf(out, "|`%s`|%s|%s|%s|\n", template.Name, template.Description, strings.Join(template.SupportedObjectKinds.ObjectKinds, ", "), params)
			}
		},
	}
)

func listCommand() *cobra.Command {
	format := common.FormatWrapper{Format: common.PlainFormat}
	c := &cobra.Command{
		Use:   "list",
		Short: "list check templates",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			knownTemplates := templates.List()
			renderFunc := formatsToRenderFuncs[format.Format]
			if renderFunc == nil {
				return errors.Errorf("unknown format: %q", format.Format)
			}
			renderFunc(knownTemplates, os.Stdout)
			return nil
		},
	}
	c.Flags().Var(&format, "format", "output format")
	return c
}

// Command defines the root of the templates command.
func Command() *cobra.Command {
	c := &cobra.Command{
		Use: "templates",
	}
	c.AddCommand(listCommand())
	return c
}
