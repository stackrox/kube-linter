package checks

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.stackrox.io/kube-linter/internal/builtinchecks"
	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/command/common"
	"golang.stackrox.io/kube-linter/internal/defaultchecks"
)

var (
	dashes = func() string {
		var sb strings.Builder
		for i := 0; i < 30; i++ {
			sb.WriteRune('-')
		}
		return sb.String()
	}()

	formatsToRenderFuncs = map[string]func([]check.Check, io.Writer) error{
		common.PlainFormat:    renderPlain,
		common.MarkdownFormat: renderMarkdown,
	}
)

func renderPlain(checks []check.Check, out io.Writer) error { //nolint:unparam // The function signature is required to match formatToRenderFuncs
	for i, chk := range checks {
		fmt.Fprintf(out, "Name: %s\nDescription: %s\nTemplate: %s\nParameters: %v\nEnabled by default: %v\n",
			chk.Name, chk.Description, chk.Template, chk.Params, defaultchecks.List.Contains(chk.Name))
		if i != len(checks)-1 {
			fmt.Fprintf(out, "\n%s\n\n", dashes)
		}
	}
	return nil
}

const (
	markDownTemplateStr = `The following table enumerates built-in checks:

| Name | Enabled by default | Description | Template | Parameters |
| ---- | ------------------ | ----------- | -------- | ---------- |
{{ range . }} | {{ .Check.Name}} | {{ if .Default }}Yes{{ else }}No{{ end }} | {{.Check.Description}} | {{.Check.Template}} | 
{{- range $key, $value := .Check.Params -}}
- {{backtick}}{{$key}}{{backtick}}: {{backtick}}{{$value}}{{backtick}} <br />
{{- else }} none {{ end -}}
|
{{ end -}}
`
)

var (
	markDownTemplate = common.MustInstantiateTemplate(markDownTemplateStr)
)

func renderMarkdown(checks []check.Check, out io.Writer) error {
	type augmentedCheck struct {
		Check   check.Check
		Default bool
	}
	augmentedChecks := make([]augmentedCheck, 0, len(checks))
	for _, chk := range checks {
		augmentedChecks = append(augmentedChecks, augmentedCheck{Check: chk, Default: defaultchecks.List.Contains(chk.Name)})
	}
	return markDownTemplate.Execute(out, augmentedChecks)
}

func listCommand() *cobra.Command {
	format := common.FormatWrapper{Format: common.PlainFormat}
	c := &cobra.Command{
		Use:   "list",
		Short: "list built-in checks",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			checks, err := builtinchecks.List()
			if err != nil {
				return err
			}
			renderFunc := formatsToRenderFuncs[format.Format]
			if renderFunc == nil {
				return errors.Errorf("unknown format: %q", format.Format)
			}
			return renderFunc(checks, os.Stdout)
		},
	}
	c.Flags().Var(&format, "format", "output format")
	return c
}

// Command defines the root of the checks command.
func Command() *cobra.Command {
	c := &cobra.Command{
		Use: "checks",
	}
	c.AddCommand(listCommand())
	return c
}
