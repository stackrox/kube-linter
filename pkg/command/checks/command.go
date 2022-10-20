package checks

import (
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.stackrox.io/kube-linter/internal/defaultchecks"
	"golang.stackrox.io/kube-linter/internal/flagutil"
	"golang.stackrox.io/kube-linter/pkg/builtinchecks"
	"golang.stackrox.io/kube-linter/pkg/command/common"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/templates"
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

{{ range . -}}
## {{ .Name}}

**Enabled by default**: {{ if isDefault . }}Yes{{ else }}No{{ end }}

**Description**: {{.Description}}

**Remediation**: {{.Remediation}}

**Template**: [{{.Template}}](templates.md#{{ templateLink . }})
{{ if .Params }}
**Parameters**:

{{ mustToYaml (default (dict) .Params ) | codeBlock "yaml" }}
{{ end -}}
{{ end -}}
`
)

var (
	checksFuncMap = template.FuncMap{
		"isDefault": func(check config.Check) bool {
			return defaultchecks.List.Contains(check.Name)
		},
		"templateLink": GetTemplateLink,
	}
	plainTemplate    = common.MustInstantiatePlainTemplate(plainTemplateStr, checksFuncMap)
	markDownTemplate = common.MustInstantiateMarkdownTemplate(markDownTemplateStr, checksFuncMap)

	formatters = common.Formatters{
		Formatters: map[common.FormatType]common.FormatFunc{
			common.PlainFormat:    plainTemplate.Execute,
			common.MarkdownFormat: markDownTemplate.Execute,
			common.JSONFormat:     common.FormatJSON,
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

// GetTemplateLink returns html anchor string for the template corresponding to the given check so that it can be used
// to reference the template section in a rendered markdown.
// E.g. template name "Deprecated Service Account Field" becomes "deprecated-service-account-field" html anchor.
func GetTemplateLink(check *config.Check) (string, error) {
	t, found := templates.Get(check.Template)
	if !found {
		return "", errors.Errorf("unexpected: check %v references non-existent template?", check)
	}
	return strings.Join(strings.Fields(strings.ToLower(t.HumanName)), "-"), nil
}
