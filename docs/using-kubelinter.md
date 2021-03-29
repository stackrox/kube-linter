# Using KubeLinter

You can run KubeLinter both locally and on your CI systems.

## Running locally

After you've [installed KubeLinter](README.md#installing-kubelinter), use the
`lint` command and provide:

## Pre-commit

KubeLinter supports the pre-commit framework, whose documentation can be found [here](https://pre-commit.com/).
This allows you to automatically call kube-linter as a git pre-commit hook.

An example `.pre-commit.config.yaml` file that uses kube-linter can be seen here:

```yaml
repos:
-   repo: https://github.com/stackrox/kube-linter
    rev: v0.1.6
    hooks:
    -   id: kube-linter
```

The currently supported pre-commit hooks are:

1. `kube-linter` which clones and builds the kube-linter utility locally using `go get`. No additional installation required.
2. `kube-linter-system` runs the kube-linter binary that must be on your systems path (for example installed through homebrew)
3. `kube-linter-docker` will pull the official docker image and run the kube-linter command that way, mounting the projects directory into
    the container. This requires docker to be installed on your machine.

<!-- tabs:start -->

#### ** Kubernetes **
- The path to your Kubernetes `yaml` file:
  ```bash
  kube-linter lint /path/to/yaml-file.yaml
  ```
- The path to a directory containing your Kubernetes `yaml` files:
  ```bash
  kube-linter lint /path/to/directory/containing/yaml-files/
  ```

#### ** Helm **
The path to a directory containing the `Chart.yaml` file:
```bash
kube-linter lint /path/to/directory/containing/Chart.yaml-file/
```

<!-- tabs:end -->

## KubeLinter commands

This section covers kube-linter command syntax, describes the command
operations, and provides some common command examples.

### Commands syntax

Use the following syntax to run KubeLinter commands:
```bash
kube-linter [resource] [command] [options]
```

where `resource`, `command`, and `options` are:

- `resource` specifies the resources on which you want to perform operations.
  For example, `checks` or `templates`.
- `command` specifies the operation that you want to perform, for example,
  `kube-linter lint`. Or the operation you want to perform on a specified
  resource, for example, `kube-linter checks list`.
- `options` specifies options for each command. For example, you can use the `-c`
  or `--config` option to specify a configuration file.

### Viewing help

Use the `--help` (or `-h`) option to get a complete list of resources and
commands. The `--help` option lists CLI reference information about available
resources, commands, and their options.

For example,

- to find all resources, run the following command:
  ```bash
  kube-linter --help
  ```
- to find available commands for a specific resource, run the following command:
  ```bash
  kube-linter checks --help
  ```
- to find available options for a specific command, run the following command:
  ```bash
  kube-linter lint --help
  ```

