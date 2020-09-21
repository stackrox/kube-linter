#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "${SCRIPT_DIR}/utils.sh"

set -eo pipefail

[[ -n "${GITHUB_TOKEN}" ]] || die "No GITHUB_TOKEN found."

tag="$1"
bin_dir="$2"
[[ -n "${tag}" && -n "${bin_dir}" ]] || die "Usage: $0 <tag> <binary_directory>"


status_code=""
curl_output=""

# gh_curl makes an authenticated request to the GitHub API at the given subpath (with the common repo prefix automatically
# added), and populates the results in the global status_code and curl_output variables.
function gh_curl() {
  local sub_path=$1
  shift
  curl -sS -w '%{http_code}' -o /tmp/curl.out -H "Authorization: Bearer ${GITHUB_TOKEN}" "https://api.github.com/repos/stackrox/kube-linter/${sub_path}" "$@" > /tmp/curl_status_code.out
  status_code="$(cat /tmp/curl_status_code.out)"
  curl_output="$(cat /tmp/curl.out)"
  rm /tmp/curl_status_code.out /tmp/curl.out
  einfo "Request to ${sub_path} with args \"$*\". Response code: ${status_code}"
}

einfo "Attempting to create release"
gh_curl releases -X POST -d '{"tag_name": "'"${tag}"'", "name": "v'"${tag}"'", "prerelease": true}'
if [[ "${status_code}" -eq 201 ]]; then
  einfo "Successfully created release."
elif [[ "${status_code}" -eq 422 ]]; then
  ewarn "Release already exists; continuing."
fi

gh_curl "releases/tags/${tag}"
if [[ "${status_code}" -ne 200 ]]; then
  die "Unexpected status code from fetching the release: ${status_code}. Output: ${curl_output}"
fi

upload_url="$(jq -r '.upload_url' <<<"${curl_output}" | sed 's/{.*$//g')"
einfo "Got upload_url: ${upload_url}"

function upload_file() {
  filename="$1"
  basename="$(basename "${filename}")"
  curl -H "Authorization: Bearer ${GITHUB_TOKEN}" "${upload_url}?name=${basename}" -X POST -F "file=@${filename}"
  echo
  rm "${filename}"
}

einfo "Uploading darwin binary..."
tar -czf /tmp/kube-linter-darwin.tar.gz "${bin_dir}/darwin/kube-linter"
upload_file /tmp/kube-linter-darwin.tar.gz

einfo "Uploading linux binary..."
tar -czf /tmp/kube-linter-linux.tar.gz "${bin_dir}/linux/kube-linter"
upload_file /tmp/kube-linter-linux.tar.gz

einfo "Uploading windows binary..."
tar -czf /tmp/kube-linter-windows.tar.gz "${bin_dir}/windows/kube-linter.exe"
upload_file /tmp/kube-linter-windows.tar.gz
