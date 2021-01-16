package lint

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
	"golang.stackrox.io/kube-linter/pkg/builtinchecks"
	"golang.stackrox.io/kube-linter/pkg/checkregistry"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/configresolver"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/run"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Command is the command for the lint command.
func Command() *cobra.Command {
	var configPath string
	var verbose bool

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
			lintCtxs, err := lintcontext.CreateContexts(args...)
			if err != nil {
				return err
			}
			if verbose {
				for _, lintCtx := range lintCtxs {
					for _, invalidObj := range lintCtx.InvalidObjects() {
						fmt.Fprintf(os.Stderr, "Warning: failed to load object from %s: %v\n", invalidObj.Metadata.FilePath, invalidObj.LoadErr)
					}
				}
			}
			var atLeastOneObjectFound bool
			for _, lintCtx := range lintCtxs {
				if len(lintCtx.Objects()) > 0 {
					atLeastOneObjectFound = true
					break
				}
			}
			if !atLeastOneObjectFound {
				fmt.Fprintln(os.Stderr, "Warning: no valid objects found.")
				return nil
			}
			result, err := run.Run(lintCtxs, checkRegistry, enabledChecks)
			if err != nil {
				return err
			}
			if len(result.Reports) == 0 {
				fmt.Fprintln(os.Stderr, "No lint errors found!")
				return nil
			}
			stderrIsTerminal := terminal.IsTerminal(int(os.Stderr.Fd()))
			for _, report := range result.Reports {
				if stderrIsTerminal {
					report.FormatToTerminal(os.Stderr)
				} else {
					report.FormatPlain(os.Stderr)
				}
			}
			return errors.Errorf("found %d lint errors", len(result.Reports))
		},
	}

	c.Flags().StringVar(&configPath, "config", "", "Path to config file")
	c.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")

	config.AddFlags(c, v)
	return c
}
