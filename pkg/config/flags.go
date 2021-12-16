// Code generated by kube-linter flag codegen. DO NOT EDIT.
// +build !flagcodegen

package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AddFlags, walks through config.Check struct and bind its Member to Cobra command 
// and add respective Viper flag 
func AddFlags(c *cobra.Command, v *viper.Viper) {
	c.Flags().Bool("add-all-built-in", false, "AddAllBuiltIn, if set, adds all built-in checks. This allows users to explicitly opt-out of checks that are not relevant using Exclude.")
	if err := v.BindPFlag("checks.addAllBuiltIn", c.Flags().Lookup("add-all-built-in")); err != nil {
		panic(err)
	}
	c.Flags().Bool("do-not-auto-add-defaults", false, "DoNotAutoAddDefaults, if set, prevents the automatic addition of default checks.")
	if err := v.BindPFlag("checks.doNotAutoAddDefaults", c.Flags().Lookup("do-not-auto-add-defaults")); err != nil {
		panic(err)
	}
	c.Flags().StringSlice("exclude", nil, "Exclude is a list of check names to exclude.")
	if err := v.BindPFlag("checks.exclude", c.Flags().Lookup("exclude")); err != nil {
		panic(err)
	}
	c.Flags().StringSlice("include", nil, "Include is a list of check names to include. If a check is in both Include and Exclude, Exclude wins.")
	if err := v.BindPFlag("checks.include", c.Flags().Lookup("include")); err != nil {
		panic(err)
	}
	c.Flags().Bool("", false, "")
	if err := v.BindPFlag("allowOpenshiftKinds", c.Flags().Lookup("allowOpenshiftKinds")); err != nil {
		panic(err)
	}
}
