package templates

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

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
	markDownTemplateStr = `This page lists supported check templates.

{{ range . -}}
## {{ .HumanName }}

**Key**: {{ backtick }}{{ .Key }}{{ backtick }}

**Description**: {{ .Description }}

**Supported Objects**: {{ join "," .SupportedObjectKinds.ObjectKinds }}

**Parameters**:
{{ backtick }}{{ backtick }}{{ backtick }}
{{ getParametersJSON .Parameters }}
{{ backtick }}{{ backtick }}{{ backtick }}

{{ end -}}
`
)

var (
	markDownTemplate = common.MustInstantiateTemplate(markDownTemplateStr, template.FuncMap{
		"getParametersJSON": func(params []check.ParameterDesc) (string, error) {
			out := make([]check.HumanReadableParamDesc, 0, len(params))
			for _, param := range params {
				out = append(out, param.HumanReadableFields())
			}
			var buf bytes.Buffer
			enc := json.NewEncoder(&buf)
			enc.SetIndent("", "\t")
			if err := enc.Encode(out); err != nil {
				return "", err
			}
			return buf.String(), nil
		},
	})
)

func renderParameters(numTabs int, params []check.ParameterDesc, out io.Writer) {
	tabs := stringutils.Repeat("\t", numTabs)
	for _, param := range params {
		fmt.Fprintf(out, "%s%s:\n%s\tDescription: %s\n%s\tRequired: %v\n", tabs, param.Name, tabs, param.Description, tabs, param.Required)
		if len(param.Examples) > 0 {
			quotedExamples := make([]string, 0, len(param.Examples))
			for _, ex := range param.Examples {
				quotedExamples = append(quotedExamples, fmt.Sprintf(`"%s"`, ex))
			}
			fmt.Fprintf(out, "%s\tExample values: %s\n", tabs, strings.Join(quotedExamples, ", "))
		}
		if len(param.SubParameters) > 0 {
			fmt.Fprintf(out, "%s\tSub-parameters:\n", tabs)
			renderParameters(numTabs+1, param.SubParameters, out)
		}
	}
}

func renderPlain(templates []check.Template, out io.Writer) error { //nolint:unparam // The function signature is required to match formatToRenderFuncs
	for i, template := range templates {
		fmt.Fprintf(out, "Name: %s\nKey: %s\nDescription: %s\nSupported Objects: %v\n", template.HumanName, template.Key, template.Description, template.SupportedObjectKinds.ObjectKinds)
		if len(template.Parameters) == 0 {
			fmt.Fprintln(out, "Parameters: none")
		} else {
			fmt.Fprintf(out, "Parameters:\n")
			renderParameters(1, template.Parameters, out)
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
