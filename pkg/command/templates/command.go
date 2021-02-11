package templates

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.stackrox.io/kube-linter/internal/flagutil"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/command/common"
	"golang.stackrox.io/kube-linter/pkg/templates"
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

const (
	markDownTemplateStr = `# KubeLinter templates

KubeLinter supports the following templates:

{{ range . -}}
## {{ .HumanName }}

**Key**: {{ .Key | codeSnippet }}

**Description**: {{ .Description }}

**Supported Objects**: {{ join "," .SupportedObjectKinds.ObjectKinds }}

**Parameters**:

{{ getParametersJSON .Parameters | codeBlock "json" }}

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
Parameters:{{ range .Parameters }}{{ template "Param" .HumanReadableFields }}{{else}} none{{end}}
{{end -}}
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
