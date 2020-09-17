# kube-linter

kube-linter is a static analysis tool that checks Kubernetes YAML files to ensure the applications represented in them
adhere to best practices.

## Install

Binary downloads of can be found on [the Releases page](https://github.com/stackrox/kube-linter/releases).

Download the `kube-linter` binary, and add it to your PATH.

## Usage

To lint directories or files, simply run `./kube-linter lint files_or_dirs ...`. If a directory is passed, all files
with `.yaml` or `.yml` extensions are parsed, and Kubernetes objects are loaded from them. If a file is passed,
it is parsed irrespective of extension.

