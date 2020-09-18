
# Documentation

Welcome to the `kube-linter` documentation. Read on for more detailed information about using and configuring the tool.

## Exporing the CLI

You can run `kube-linter --help` to see a list of supported commands and flags. For each subcommand, you can
run `kube-linter <subcommand> --help` to see detailed help text and flags for it.

## Running the linter

To lint directories or files, simply run `./kube-linter lint files_or_dirs ...`. If a directory is passed, all files
with `.yaml` or `.yml` extensions are parsed, and Kubernetes objects are loaded from them. If a file is passed,
it is parsed irrespective of extension.

Users can pass a config file using the `--config` file to control which checks are executed, and to configure custom checks.
An example config file is provided [here](../config.yaml.example).

## Built-in checks 

`kube-linter` comes with a list of built-in checks, which you can find [here](generated/checks.md). Only some
built-in checks are enabled by default -- others must be explicitly enabled in the config.

## Custom checks

### Check Templates

In `kube-linter`, checks are concrete realizations of check templates. A check template describes a class of check -- it
contains logic (written in Go code) that would execute the check, and lays out (zero or more) parameters that it takes.

The list of supported check templates, along with their metadata, can be found [here](generated/templates.md).

### Custom checks

All checks in `kube-linter` are defined by referencing a check template, passing parameters to it, and adding additional
check specific metadata (like check name and description). Users can configure custom checks the same way built-in checks
are configured, and add them to the config file. The built-in checks are specified in [internal/builtinchecks](internal/builtinchecks). 
