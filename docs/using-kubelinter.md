# Using KubeLinter

You can run KubeLinter both locally and on your CI systems.

## Running locally

After you've [installed KubeLinter](README.md#installing-kubelinter), use the
`lint` command and provide:

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


> [!NOTE] To get structured output, use the `--format` option.
> For example,
> - Use `--format=json` to get the output in JSON format.
> - Use `--format=sarif` to get the output in the [SARIF spec](https://github.com/microsoft/sarif-tutorials).

## Using KubeLinter with the pre-commit framework

If you are using the [pre-commit framework](https://pre-commit.com/) for
managing Git pre-commit hooks, you can install and use KubeLinter as a
pre-commit hook. To do this, add the following to your `.pre-commit-config.yaml`:

```yaml
  - repo: https://github.com/stackrox/kube-linter
    rev: 0.6.0 # kube-linter version 
    hooks:
        # You can change this to kube-linter-system or kube-linter-docker
      - id: kube-linter
```

The [`.pre-commit-hooks.yaml`](https://raw.githubusercontent.com/stackrox/kube-linter/main/.pre-commit-hooks.yaml)
includes the following pre-commit hooks:

1. `kube-linter`: Clones, builds, and installs KubeLinter locally (by using `go
   get`).
2. `kube-linter-system`: Runs the KubeLinter binary.
   - You must [install KubeLinter](README.md#installing-kubelinter) and add it
     to your path before running this pre-commit hook.
3. `kube-linter-docker`: Pulls the KubeLinter docker image, mounts the project
   directory in the container, and runs the `kube-linter` command.
   - You must [install Docker](https://docs.docker.com/engine/install/) before
     running this pre-commit hook.

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

