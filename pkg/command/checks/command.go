package checks

import (
	"os"
	"sort"
	"text/template"

	"github.com/spf13/cobra"
	"golang.stackrox.io/kube-linter/internal/defaultchecks"
	"golang.stackrox.io/kube-linter/internal/flagutil"
	"golang.stackrox.io/kube-linter/pkg/builtinchecks"
	"golang.stackrox.io/kube-linter/pkg/command/common"
	"golang.stackrox.io/kube-linter/pkg/config"
)

const (
	plainTemplateStr = `{{ range $i, $_ := . }}
{{- if $i}}
------------------------------

{{end -}}
Name: {{.Name}}
Description: {{.Description}}
Remediation: {{.Remediation}}
Template: {{.Template}}
Parameters: {{.Params}}
Enabled by default: {{ isDefault . }}
{{end -}}
`

	markDownTemplateStr = `# KubeLinter checks

KubeLinter includes the following built-in checks:

| Name | Enabled by default | Description | Remediation | Template | Parameters |
| ---- | ------------------ | ----------- | ----------- | -------- | ---------- |
{{ range . }} | {{.Name}} | {{ if isDefault . }}Yes{{ else }}No{{ end }} | {{.Description}} | {{.Remediation}} | {{.Template}} | {{ mustToJson (default (dict) .Params ) | codeSnippetInTable }} |
{{ end -}}
`
)

var (
	checksFuncMap = template.FuncMap{
		"isDefault": func(check config.Check) bool {
			return defaultchecks.List.Contains(check.Name)
		},
	}
	plainTemplate    = common.MustInstantiatePlainTemplate(plainTemplateStr, checksFuncMap)
	markDownTemplate = common.MustInstantiateMarkdownTemplate(markDownTemplateStr, checksFuncMap)

	formatters = common.Formatters{
		Formatters: map[common.FormatType]common.FormatFunc{
			common.PlainFormat:    plainTemplate.Execute,
			common.MarkdownFormat: markDownTemplate.Execute,
			common.JsonFormat:     common.FormatJson,
		},
	}
)

func listCommand() *cobra.Command {
	format := flagutil.NewEnumFlag("Output format", formatters.GetEnabledFormatters(), common.PlainFormat)
	c := &cobra.Command{
		Use:   "list",
		Short: "List built-in checks",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			checks, err := builtinchecks.List()
			if err != nil {
				return err
			}
			sort.Slice(checks, func(i, j int) bool {
				return checks[i].Name < checks[j].Name
			})
			renderFunc, err := formatters.FormatterByType(format.String())
			if err != nil {
				return err
			}
			return renderFunc(os.Stdout, checks)
		},
	}
	c.Flags().Var(format, "format", format.Usage())
	return c
}

// Command defines the root of the checks command.
func Command() *cobra.Command {
	c := &cobra.Command{
		Use:   "checks",
		Short: "View more information on lint checks",
	}
	c.AddCommand(listCommand())
	return c
}
