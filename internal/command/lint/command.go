package lint

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.stackrox.io/kube-linter/internal/builtinchecks"
	"golang.stackrox.io/kube-linter/internal/checkregistry"
	"golang.stackrox.io/kube-linter/internal/config"
	"golang.stackrox.io/kube-linter/internal/configresolver"
	"golang.stackrox.io/kube-linter/internal/lintcontext"
	"golang.stackrox.io/kube-linter/internal/run"
	"golang.stackrox.io/kube-linter/internal/utils"
)

// Command is the command for the lint command.
func Command() *cobra.Command {
	var dir string
	var configPath string
	c := &cobra.Command{
		Use:  "lint",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
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
			if err := lintCtx.LoadObjectsFromDir(dir); err != nil {
				return err
			}
			if len(lintCtx.Objects) == 0 {
				fmt.Fprintln(os.Stderr, "Warning: no objects found.")
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
			for _, report := range result.Reports {
				report.FormatTo(os.Stderr)
			}
			os.Exit(1)
			return nil
		},
	}
	c.Flags().StringVar(&dir, "dir", "", "directory of YAML files to lint")
	utils.Must(
		c.MarkFlagRequired("dir"),
	)
	c.Flags().StringVar(&configPath, "config", "", "path to config file")
	return c
}
