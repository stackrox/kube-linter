#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "${SCRIPT_DIR}/utils.sh"

set -eo pipefail

[[ -n "${GITHUB_TOKEN}" ]] || die "No GITHUB_TOKEN found."

tag="$1"
bin_dir="$2"
[[ -n "${tag}" && -n "${bin_dir}" ]] || die "Usage: $0 <tag> <binary_directory>"

mkdir /tmp/release-artifacts
tar -czf /tmp/release-artifacts/kube-linter-darwin.tar.gz "${bin_dir}/darwin/kube-linter"
tar -czf /tmp/release-artifacts/kube-linter-linux.tar.gz "${bin_dir}/linux/kube-linter"
tar -czf /tmp/release-artifacts/kube-linter-windows.tar.gz "${bin_dir}/windows/kube-linter.exe"

ghr -prerelease -n "v${tag}" "${tag}" /tmp/release-artifacts
rm -r /tmp/release-artifacts
