package checks

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.stackrox.io/kube-linter/internal/defaultchecks"
	"golang.stackrox.io/kube-linter/internal/stringutils"
	"golang.stackrox.io/kube-linter/pkg/builtinchecks"
	"golang.stackrox.io/kube-linter/pkg/command/common"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/templates"
)

var (
	dashes = stringutils.Repeat("-", 30)

	formatsToRenderFuncs = map[string]func([]config.Check, io.Writer) error{
		common.PlainFormat:    renderPlain,
		common.MarkdownFormat: renderMarkdown,
	}
)

func renderPlain(checks []config.Check, out io.Writer) error { //nolint:unparam // The function signature is required to match formatToRenderFuncs
	for i, chk := range checks {
		fmt.Fprintf(out, "Name: %s\nDescription: %s\nRemediation: %s\nTemplate: %s\nParameters: %v\nEnabled by default: %v\n",
			chk.Name, chk.Description, chk.Remediation, chk.Template, chk.Params, defaultchecks.List.Contains(chk.Name))
		if i != len(checks)-1 {
			fmt.Fprintf(out, "\n%s\n\n", dashes)
		}
	}
	return nil
}

const (
	markDownTemplateStr = `# KubeLinter checks

KubeLinter includes the following built-in checks:

{{ range . -}}
## {{ .Check.Name}}

**Enabled by default**: {{ if .Default }}Yes{{ else }}No{{ end }}

**Description**: {{.Check.Description}}

**Remediation**: {{.Check.Remediation}}

**Template**: [{{.Check.Template}}](generated/templates.md#{{.TemplateLink}})

**Parameters**:
{{ mustToJson (default (dict) .Check.Params ) | codeBlock }}

{{ end -}}
`
)

var (
	markDownTemplate = common.MustInstantiateTemplate(markDownTemplateStr, nil)
)

func renderMarkdown(checks []config.Check, out io.Writer) error {
	type augmentedCheck struct {
		Check        config.Check
		Default      bool
		TemplateLink string
	}
	augmentedChecks := make([]augmentedCheck, 0, len(checks))
	for _, chk := range checks {
		template, found := templates.Get(chk.Template)
		if !found {
			return errors.Errorf("unexpected: check %v references non-existent template?", chk)
		}
		augmentedChecks = append(augmentedChecks, augmentedCheck{
			Check:        chk,
			Default:      defaultchecks.List.Contains(chk.Name),
			TemplateLink: strings.Join(strings.Fields(strings.ToLower(template.HumanName)), "-"),
		})
	}
	return markDownTemplate.Execute(out, augmentedChecks)
}

func listCommand() *cobra.Command {
	format := common.FormatValueFactory(common.PlainFormat)
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
			renderFunc := formatsToRenderFuncs[format.String()]
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
