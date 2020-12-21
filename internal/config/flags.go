// Code generated by kube-linter flag codegen. DO NOT EDIT.
// +build !flagcodegen

package config

import (
	"github.com/spf13/viper"
	"github.com/spf13/cobra"
)

// AddFlags, walks through config.Check struct and bind its Member to Cobra command 
// and add respective Viper flag 
func AddFlags(c *cobra.Command, v *viper.Viper) {
		c.Flags().Bool("add-all-built-in", false, "AddAllBuiltIn, if set, adds all built-in checks. This allows users to explicitly opt-out of checks that are not relevant using Exclude.")
		v.BindPFlag("checks.add-all-built-in", c.Flags().Lookup("add-all-built-in"))
		c.Flags().Bool("do-not-auto-add-defaults", false, "DoNotAutoAddDefaults, if set, prevents the automatic addition of default checks.")
		v.BindPFlag("checks.do-not-auto-add-defaults", c.Flags().Lookup("do-not-auto-add-defaults"))		
		c.Flags().StringSlice("exclude", nil, "Exclude is a list of check names to exclude.")
		v.BindPFlag("checks.exclude", c.Flags().Lookup("exclude"))		
		c.Flags().StringSlice("include", nil, "Include is a list of check names to include. If a check is in both Include and Exclude, Exclude wins.")
		v.BindPFlag("checks.include", c.Flags().Lookup("include"))
}
