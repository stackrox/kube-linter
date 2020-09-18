# kube-linter

kube-linter is a static analysis tool that checks Kubernetes YAML files to ensure the applications represented in them
adhere to best practices.

In detail, `kube-linter` is a binary that takes in paths to YAML files, and runs a list of checks
against them. If any lint errors are found, they are printed to standard error, and `kube-linter` returns a non-zero 
exit code.

The list of checks that is run is configurable. `kube-linter` comes with several built-in checks, only some of which
are enabled by default. Users can also create custom checks.

## Install

If you have `go` installed, you can run `go get golang.stackrox.io/kube-linter/cmd/kube-linter`.

`kube-linter` binaries can be downloaded from [the Releases page](https://github.com/stackrox/kube-linter/releases).
Download the `kube-linter` binary, and add it to your PATH.

## Usage

To lint directories or files, simply run `./kube-linter lint files_or_dirs ...`. If a directory is passed, all files
with `.yaml` or `.yml` extensions are parsed, and Kubernetes objects are loaded from them. If a file is passed,
it is parsed irrespective of extension.

Users can pass a config file using the `--config` file to control which checks are executed, and to configure custom checks.
An example config file is provided in `config.yaml.example`.

See the [documentation](./docs) for more details.

# WARNING: Breaking changes possible

kube-linter is currently in a very early stage of development. There may be breaking changes to the command usage, flags
and config file formats.