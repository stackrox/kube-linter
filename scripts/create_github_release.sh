#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "${SCRIPT_DIR}/utils.sh"

set -eo pipefail

[[ -n "${GITHUB_TOKEN}" ]] || die "No GITHUB_TOKEN found."

tag="$1"
bin_dir="$2"
[[ -n "${tag}" && -n "${bin_dir}" ]] || die "Usage: $0 <tag> <binary_directory>"

tmp_dir="$(mktemp -d)"
for os in darwin linux windows; do
  bin_name="kube-linter"
  if [[ "${os}" == "windows" ]]; then
    bin_name="kube-linter.exe"
  fi
  tar -C "${bin_dir}/${os}" -czf "${tmp_dir}/kube-linter-${os}.tar.gz" "${bin_name}"
  zip "${tmp_dir}/kube-linter-${os}.zip" "${bin_dir}/${os}/${bin_name}"
done

ghr -prerelease -n "v${tag}" "${tag}" "${tmp_dir}"
rm -rf "${tmp_dir}"
