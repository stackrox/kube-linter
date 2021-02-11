package templates

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.stackrox.io/kube-linter/internal/flagutil"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/command/common"
	"golang.stackrox.io/kube-linter/pkg/templates"
)

const (
	markDownTemplateStr = `# KubeLinter templates

KubeLinter supports the following templates:

{{ range . -}}
## {{ .HumanName }}

**Key**: {{ .Key | codeSnippet }}

**Description**: {{ .Description }}

**Supported Objects**: {{ join "," .SupportedObjectKinds.ObjectKinds }}

**Parameters**:

{{ toPrettyJson .HumanReadableParameters | codeBlock "json" }}

{{ end -}}
`

	plainTemplateStr = `{{- define "Param" }}{{ $tabs := repeat .NestingLevel "\t" }}
	{{$tabs}}{{.Name}}:
		{{$tabs}}Description: {{.Description}}
		{{$tabs}}Required: {{.Required}}{{if .Examples}}
		{{$tabs}}Example values: {{ range $i, $_ := .Examples }}{{if $i}}, {{end}}{{ printf "%q" . }}{{end}}{{end}}{{if .SubParameters}}
		{{$tabs}}Sub-parameters: {{ range .SubParameters }}{{ template "Param" . }}{{end}}{{end}}{{if .ArrayElemType}}
		{{$tabs}}Elem type: {{.ArrayElemType}}{{end}}
{{- end -}}
{{ range $i, $_ := . }}
{{- if $i}}
------------------------------

{{end -}}
Name: {{.HumanName}}
Key: {{.Key}}
Description: {{.Description}}
Supported Objects: {{.SupportedObjectKinds.ObjectKinds}}
Parameters:{{ range .HumanReadableParameters }}{{ template "Param" . }}{{else}} none{{end}}
{{end -}}
`
)

var (
	outputFormats = flagutil.NewEnumValueFactory("Output format", []string{common.PlainFormat, common.MarkdownFormat, common.JsonFormat})

	formatters = map[string]func([]check.Template, io.Writer) error{
		common.PlainFormat:    renderPlain,
		common.MarkdownFormat: renderMarkdown,
		common.JsonFormat: func(templates []check.Template, out io.Writer) error {
			return common.FormatJson(templates, out)
		},
	}
)

var (
	markDownTemplate = common.MustInstantiateTemplate(markDownTemplateStr, nil)

	plainTemplate = common.MustInstantiateTemplate(plainTemplateStr, nil)
)

func renderPlain(templates []check.Template, out io.Writer) error {
	return plainTemplate.Execute(out, templates)
}

func renderMarkdown(templates []check.Template, out io.Writer) error {
	return markDownTemplate.Execute(out, templates)
}

func listCommand() *cobra.Command {
	format := outputFormats(common.PlainFormat)
	c := &cobra.Command{
		Use:   "list",
		Short: "List check templates",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			knownTemplates := templates.List()
			renderFunc := formatters[format.String()]
			if renderFunc == nil {
				return errors.Errorf("unknown format: %q", format.String())
			}
			return renderFunc(knownTemplates, os.Stdout)
		},
	}
	c.Flags().Var(format, "format", format.Usage())
	return c
}

// Command defines the root of the templates command.
func Command() *cobra.Command {
	c := &cobra.Command{
		Use:   "templates",
		Short: "View more information on check templates",
	}
	c.AddCommand(listCommand())
	return c
}
