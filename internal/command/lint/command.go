package lint

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.stackrox.io/kube-linter/internal/config"
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
			cfg, err := config.Load(configPath)
			if err != nil {
				return errors.Wrap(err, "failed to load config")
			}
			lintCtx := lintcontext.New()
			if err := lintCtx.LoadObjectsFromDir(dir); err != nil {
				return err
			}
			result, err := run.Run(lintCtx, cfg)
			if err != nil {
				return err
			}
			if len(result.Reports) == 0 {
				fmt.Fprintln(os.Stderr, "No lint errors found!")
				return nil
			}
			for _, report := range result.Reports {
				fmt.Fprintln(os.Stderr, report.Format())
			}
			os.Exit(1)
			return nil
		},
	}
	c.Flags().StringVar(&dir, "dir", "", "directory of YAML files to lint")
	c.Flags().StringVar(&configPath, "config", "", "path to config file")
	utils.Must(
		c.MarkFlagRequired("dir"),
		c.MarkFlagRequired("config"),
	)
	return c
}
