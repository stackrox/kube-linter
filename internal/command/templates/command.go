package templates

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/command/common"
	"golang.stackrox.io/kube-linter/internal/stringutils"
	"golang.stackrox.io/kube-linter/internal/templates"
)

var (
	dashes = stringutils.Repeat("-", 30)

	formatsToRenderFuncs = map[string]func([]check.Template, io.Writer) error{
		common.PlainFormat:    renderPlain,
		common.MarkdownFormat: renderMarkdown,
	}
)

const (
	markDownTemplateStr = `The following table enumerates supported check templates:

| Name | Description | Supported Objects | Parameters |
| ---- | ----------- | ----------------- | ---------- |
{{ range . }} | {{ .Name}} | {{ .Description }} | {{ join "," .SupportedObjectKinds.ObjectKinds }} |
{{- range .Parameters -}}
- {{backtick}}{{.ParamName}}{{backtick}}{{ if .Required }} (required){{ end }}: {{ .Description }} <br />
{{- else }} none {{ end -}}
|
{{ end -}}
`
)

var (
	markDownTemplate = common.MustInstantiateTemplate(markDownTemplateStr)
)

func renderPlain(templates []check.Template, out io.Writer) error { //nolint:unparam // The function signature is required to match formatToRenderFuncs
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
	return nil
}

func renderMarkdown(templates []check.Template, out io.Writer) error {
	return markDownTemplate.Execute(out, templates)
}

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
			return renderFunc(knownTemplates, os.Stdout)
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
