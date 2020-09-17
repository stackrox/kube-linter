package lint

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"golang.stackrox.io/kube-linter/internal/builtinchecks"
	"golang.stackrox.io/kube-linter/internal/checkregistry"
	"golang.stackrox.io/kube-linter/internal/config"
	"golang.stackrox.io/kube-linter/internal/configresolver"
	"golang.stackrox.io/kube-linter/internal/lintcontext"
	"golang.stackrox.io/kube-linter/internal/run"
)

// Command is the command for the lint command.
func Command() *cobra.Command {
	var configPath string
	var verbose bool

	c := &cobra.Command{
		Use:  "lint",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			checkRegistry := checkregistry.New()
			if err := builtinchecks.LoadInto(checkRegistry); err != nil {
				return err
			}
			var cfg config.Config
			if configPath != "" {
				var err error
				cfg, err = config.Load(configPath)
				if err != nil {
					return errors.Wrap(err, "failed to load config")
				}
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
			lintCtx := lintcontext.New()
			for _, dir := range args {
				if err := lintCtx.LoadObjectsFromPath(dir); err != nil {
					return err
				}
			}
			if verbose {
				for _, invalidObj := range lintCtx.InvalidObjects() {
					fmt.Fprintf(os.Stderr, "Warning: failed to load object from %s: %v\n", invalidObj.Metadata.FilePath, invalidObj.LoadErr)
				}
			}
			if len(lintCtx.Objects()) == 0 {
				fmt.Fprintln(os.Stderr, "Warning: no valid objects found.")
				return nil
			}
			result, err := run.Run(lintCtx, checkRegistry, enabledChecks)
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
	c.Flags().StringVar(&configPath, "config", "", "path to config file")
	c.Flags().BoolVarP(&verbose, "verbose", "v", false, "whether to enable verbose logs")
	return c
}
