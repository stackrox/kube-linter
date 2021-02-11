package templates

import (
	"os"

	"github.com/spf13/cobra"
	"golang.stackrox.io/kube-linter/internal/flagutil"
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
	markDownTemplate = common.MustInstantiateTemplate(markDownTemplateStr, nil)
	plainTemplate    = common.MustInstantiateTemplate(plainTemplateStr, nil)

	formatters = common.Formatters{
		Formatters: map[common.FormatType]common.FormatFunc{
			common.PlainFormat:    plainTemplate.Execute,
			common.MarkdownFormat: markDownTemplate.Execute,
			common.JsonFormat:     common.FormatJson,
		},
	}

	outputFormats = flagutil.NewEnumValueFactory("Output format", formatters.GetEnabledFormatters())
)

func listCommand() *cobra.Command {
	format := outputFormats(common.PlainFormat)
	c := &cobra.Command{
		Use:   "list",
		Short: "List check templates",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			knownTemplates := templates.List()
			formatFunc, err := formatters.FormatterByType(format.String())
			if err != nil {
				return err
			}
			return formatFunc(os.Stdout, knownTemplates)
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
