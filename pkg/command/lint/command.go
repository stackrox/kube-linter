package lint

import (
	"fmt"
	"os"

	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/pathutil"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"golang.stackrox.io/kube-linter/internal/flagutil"
	"golang.stackrox.io/kube-linter/pkg/builtinchecks"
	"golang.stackrox.io/kube-linter/pkg/checkregistry"
	"golang.stackrox.io/kube-linter/pkg/command/common"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/configresolver"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/run"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	plainTemplateStr = `KubeLinter {{.Summary.KubeLinterVersion}}

{{range .Reports}}
{{- .Object.Metadata.FilePath | bold}}: (object: {{.Object.GetK8sObjectName | bold}}) {{.Diagnostic.Message | red}} (check: {{.Check | yellow}}, remediation: {{.Remediation | yellow}})

{{else}}No lint errors found!
{{end -}}
`
)

var (
	plainTemplate = common.MustInstantiatePlainTemplate(plainTemplateStr, nil)

	formatters = common.Formatters{
		Formatters: map[common.FormatType]common.FormatFunc{
			common.JSONFormat:  common.FormatJSON,
			common.SARIFFormat: formatLintSarif,
			common.PlainFormat: plainTemplate.Execute,
		},
	}
)

// Command is the command for the lint command.
func Command() *cobra.Command {
	var configPath string
	var failIfNoObjects bool
	var verbose bool
	var errorOnInvalidResource bool
	format := flagutil.NewEnumFlag("Output format", formatters.GetEnabledFormatters(), common.PlainFormat)

	v := viper.New()

	c := &cobra.Command{
		Use:   "lint",
		Args:  cobra.MinimumNArgs(1),
		Short: "Lint Kubernetes YAML files and Helm charts",
		RunE: func(cmd *cobra.Command, args []string) error {
			checkRegistry := checkregistry.New()
			if err := builtinchecks.LoadInto(checkRegistry); err != nil {
				return err
			}

			// Load Configuration
			cfg, err := config.Load(v, configPath)
			if err != nil {
				return errors.Wrap(err, "failed to load config")
			}

			if err := configresolver.LoadCustomChecksInto(&cfg, checkRegistry); err != nil {
				return err
			}
			enabledChecks, err := configresolver.GetEnabledChecksAndValidate(&cfg, checkRegistry)
			if err != nil {
				return err
			}
			if len(enabledChecks) == 0 {
				fmt.Fprintln(os.Stderr, "Warning: no checks enabled.")
				return nil
			}
			ignorePaths, err := configresolver.GetIgnorePaths(&cfg)
			if err != nil {
				return err
			}

			absArgs := make([]string, 0, len(args))
			for _, arg := range args {
				if arg == lintcontext.ReadFromStdin {
					absArgs = append(absArgs, lintcontext.ReadFromStdin)
					continue
				}
				absArg, err := pathutil.GetAbsolutPath(arg)
				if err != nil {
					return err
				}
				absArgs = append(absArgs, absArg)
			}

			lintCtxs, err := lintcontext.CreateContexts(ignorePaths, absArgs...)
			if err != nil {
				return err
			}
			invalidObjectsResult := generateReportFromInvalidObjects(lintCtxs)
			if verbose {
				for _, invalidObj := range invalidObjectsResult {
					_, _ = fmt.Fprintf(os.Stderr, "Warning: failed to load object from %s: %v\n", invalidObj.Object.Metadata.FilePath, invalidObj.Diagnostic.Message)
				}
			}

			var atLeastOneObjectFound = errorOnInvalidResource && len(invalidObjectsResult) > 0
			for _, lintCtx := range lintCtxs {
				if len(lintCtx.Objects()) > 0 {
					atLeastOneObjectFound = true
					break
				}
			}

			if !atLeastOneObjectFound {
				msg := "no valid objects found"
				if failIfNoObjects {
					return errors.New(msg)
				}
				fmt.Fprintf(os.Stderr, "Warning: %s.\n", msg)
				return nil
			}
			result, err := run.Run(lintCtxs, checkRegistry, enabledChecks)
			if err != nil {
				return err
			}

			if errorOnInvalidResource {
				result.Reports = append(result.Reports, invalidObjectsResult...)
			}

			formatter, err := formatters.FormatterByType(format.String())
			if err != nil {
				return err
			}
			err = formatter(os.Stdout, result)
			if err != nil {
				return errors.Wrap(err, "output formatting failed")
			}

			if len(result.Reports) > 0 {
				err = errors.Errorf("found %d lint errors", len(result.Reports))
			}
			return err
		},
	}

	c.Flags().StringVar(&configPath, "config", "", "Path to config file")
	c.Flags().BoolVarP(&failIfNoObjects, "fail-if-no-objects-found", "", false, "Return non-zero exit code if no valid objects are found or failed to parse")
	c.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	c.Flags().Var(format, "format", format.Usage())
	c.Flags().BoolVarP(&errorOnInvalidResource, "fail-on-invalid-resource", "", false, "Error out when we have an invalid resource")

	config.AddFlags(c, v)
	return c
}

func generateReportFromInvalidObjects(lintCtxs []lintcontext.LintContext) []diagnostic.WithContext {
	var invalidObjectsResult []diagnostic.WithContext
	for _, lintCtx := range lintCtxs {
		for _, invalidObj := range lintCtx.InvalidObjects() {
			invalidObjectsResult = append(invalidObjectsResult, diagnostic.WithContext{
				Diagnostic: diagnostic.Diagnostic{
					Message: invalidObj.LoadErr.Error(),
				},
				Check:       "failed-to-load-object",
				Remediation: "Confirm that the file is accessible and is valid k8s yaml.",
				Object: lintcontext.Object{
					Metadata:  invalidObj.Metadata,
					K8sObject: &unstructured.Unstructured{},
				},
			})
		}
	}
	return invalidObjectsResult
}
