package main

import (
	"fmt"
	"os"

	"golang.stackrox.io/kube-linter/pkg/command/root"
	// Register templates
	_ "golang.stackrox.io/kube-linter/pkg/templates/all"
)

func main() {
	c := root.Command()
	if err := c.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
