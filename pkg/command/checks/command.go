package checks

import (
	"io"
	"os"
	"sort"
	"text/template"

	"github.com/pkg/errors"
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
	outputFormats = flagutil.NewEnumValueFactory("Output format", []string{common.PlainFormat, common.MarkdownFormat, common.JsonFormat})

	formatters = map[string]func([]config.Check, io.Writer) error{
		common.PlainFormat:    renderPlain,
		common.MarkdownFormat: renderMarkdown,
		common.JsonFormat: func(checks []config.Check, out io.Writer) error {
			return common.FormatJson(out, checks)
		},
	}
)

var (
	checksFuncMap = template.FuncMap{
		"isDefault": func(check config.Check) bool {
			return defaultchecks.List.Contains(check.Name)
		},
	}

	plainTemplate    = common.MustInstantiateTemplate(plainTemplateStr, checksFuncMap)
	markDownTemplate = common.MustInstantiateTemplate(markDownTemplateStr, checksFuncMap)
)

func renderPlain(checks []config.Check, out io.Writer) error {
	return plainTemplate.Execute(out, checks)
}

func renderMarkdown(checks []config.Check, out io.Writer) error {
	return markDownTemplate.Execute(out, checks)
}

func listCommand() *cobra.Command {
	format := outputFormats(common.PlainFormat)
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
			renderFunc := formatters[format.String()]
			if renderFunc == nil {
				return errors.Errorf("unknown format: %q", format.String())
			}
			return renderFunc(checks, os.Stdout)
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
