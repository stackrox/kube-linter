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

#### ** Kustomize **
The path to a directory containing a `kustomization.yaml` file:
```bash
kube-linter lint /path/to/directory/containing/kustomization.yaml
```

> [!NOTE] KubeLinter automatically detects Kustomize directories and renders the manifests before linting.
> Source file paths are preserved in lint reports, pointing to the actual base/overlay files rather than generated output.

<!-- tabs:end -->


> [!NOTE] To get structured output, use the `--format` option.
> For example,
> - Use `--format=json` to get the output in JSON format.
> - Use `--format=sarif` to get the output in the [SARIF spec](https://github.com/microsoft/sarif-tutorials).

## Multiple Output Formats

KubeLinter supports writing results in multiple formats in a single run. This eliminates the need to run the linter multiple times for different output formats, improving efficiency and ensuring consistency across reports.

### Single format to stdout (backward compatible)

The default behavior outputs a single format to stdout:

```bash
kube-linter lint --format json myapp.yaml
```

### Single format to file

To write output to a file, use the `--output` flag:

```bash
kube-linter lint --format sarif --output results.sarif myapp.yaml
```

### Multiple formats to multiple files

You can generate multiple output formats in a single run by repeating the `--format` and `--output` flags. The flags are paired by position (first `--format` with first `--output`, etc.):

```bash
kube-linter lint \
  --format sarif --output kube-linter.sarif \
  --format json --output kube-linter.json \
  --format plain --output kube-linter.txt \
  myapp.yaml
```

This command will:
- Generate a SARIF format report in `kube-linter.sarif`
- Generate a JSON format report in `kube-linter.json`
- Generate a plain text report in `kube-linter.txt`
- Process the files only once, improving efficiency

### Important Notes

- **Positional pairing**: Format and output flags are paired by position. The first `--format` corresponds to the first `--output`, the second `--format` to the second `--output`, and so on.
- **All stdout or all files**: Either all formats write to stdout (when no `--output` flags are provided), or each format must have a corresponding output file. You cannot mix stdout and file outputs.
- **Error handling**: If one format fails to write, the other successful formats are still written, and errors are reported at the end.
- **File overwrites**: Output files are created or overwritten if they already exist.
- **Duplicate formats allowed**: You can specify the same format multiple times with different output files if needed.

### Examples

**Example 1: Generate JSON and SARIF for CI integration**
```bash
kube-linter lint \
  --format json --output build/kube-linter.json \
  --format sarif --output build/kube-linter.sarif \
  --config .kube-linter.yaml \
  deployments/
```

**Example 2: Error case - mismatched counts**
```bash
# This will fail with an error
kube-linter lint \
  --format json --format sarif \
  --output out.json \
  pod.yaml
```
Output: `Error: format/output mismatch: 2 format(s) specified but 1 output(s) provided`

**Example 3: Error case - multiple formats to stdout**
```bash
# This will fail with an error (prevents unparseable mixed output)
kube-linter lint --format json --format plain pod.yaml
```
Output: `Error: multiple formats require explicit --output flags. Use --output to specify files, or use a single --format for stdout`

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

## Using KubeLinter as part of your CI pipeline

If you're using GitHub, there's the [KubeLinter GitHub Action](README.md#kubelinter-github-action) available.

Alternatively, you can grab the latest binary release with the following commands
```bash
LOCATION=$(curl -s https://api.github.com/repos/stackrox/kube-linter/releases/latest \
| jq -r '.tag_name
    | "https://github.com/stackrox/kube-linter/releases/download/\(.)/kube-linter-linux.tar.gz"')
curl -L -o kube-linter-linux.tar.gz $LOCATION
mkdir kube-linter/
tar -xf kube-linter-linux.tar.gz -C "kube-linter/"

```
and pass `--fail-on-invalid-resource` as an option to have your pipeline fail if your YAML file can't be parsed. See the following example:
```bash
./kube-linter/kube-linter lint --fail-on-invalid-resource /path/to/yaml-file.yaml
```

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

