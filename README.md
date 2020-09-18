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

Alternatively, `kube-linter` binaries can be downloaded from [the Releases page](https://github.com/stackrox/kube-linter/releases).
Download the `kube-linter` binary, and add it to your PATH.

You can also build `kube-linter` from source by cloning the repo, and running `make build`. This will compile a `kube-linter`
binary into the `bin` folder inside the repo.

Note that you will need to have the `go` command installed for this to work.
To install the Go command, follow the instructions [here](https://golang.org/doc/install).

## Usage

See the [documentation](./docs) for details on how to get started.

# WARNING: Breaking changes possible

kube-linter is currently in a very early stage of development. There may be breaking changes to the command usage, flags
and config file formats.