package templates

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.stackrox.io/kube-linter/internal/check"
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
		"plain": func(templates []check.Template, out io.Writer) {
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
		"markdown": func(templates []check.Template, out io.Writer) {
			fmt.Fprintf(out, "\n\n| Name | Description | Supported Objects | Parameters |\n --- | --- | --- | --- |\n")
			for _, template := range templates {
				var params string
				if len(template.Parameters) == 0 {
					params = "none"
				} else {
					var sb strings.Builder
					for _, param := range template.Parameters {
						sb.WriteString(fmt.Sprintf("- %s%s: %s <br />", param.ParamName, ternary.String(param.Required, " (required)", ""), param.Description))
					}
					params = sb.String()
				}
				fmt.Fprintf(out, "|%s|%s|%s|%s|\n", template.Name, template.Description, strings.Join(template.SupportedObjectKinds.ObjectKinds, ", "), params)
			}
		},
	}

	allValidFormats = func() []string {
		out := make([]string, 0, len(formatsToRenderFuncs))
		for format := range formatsToRenderFuncs {
			out = append(out, format)
		}
		return out
	}()
)

type formatWrapper struct {
	format string
}

func (f *formatWrapper) String() string {
	return f.format
}

func (f *formatWrapper) Set(input string) error {
	if _, ok := formatsToRenderFuncs[input]; !ok {
		return errors.Errorf("%q is not a valid option (valid options are %v)", input, allValidFormats)
	}
	f.format = input
	return nil
}

func (f *formatWrapper) Type() string {
	return "output format"
}

func listCommand() *cobra.Command {
	format := formatWrapper{format: "plain"}
	c := &cobra.Command{
		Use:  "list",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			knownTemplates := templates.List()
			renderFunc := formatsToRenderFuncs[format.format]
			if renderFunc == nil {
				return errors.Errorf("unknown format: %q", format.format)
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
