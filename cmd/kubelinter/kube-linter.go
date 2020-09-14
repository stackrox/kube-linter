package main

import (
	"os"

	"golang.stackrox.io/kube-linter/internal/command/root"
	// Register templates
	_ "golang.stackrox.io/kube-linter/internal/templates/all"
)

func main() {
	c := root.Command()
	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}
