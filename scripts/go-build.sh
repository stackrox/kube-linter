#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "${SCRIPT_DIR}/utils.sh"

set -euo pipefail

main_srcdir="$1"
[[ -n "${main_srcdir}" ]] || die "Usage: $0 <directory with main file>"

x_defs=()
x_def_errors=()

while read -r line || [[ -n "$line" ]]; do
  if [[ "$line" =~ ^[[:space:]]*$ ]]; then
    continue
  elif [[ "$line" =~ ^([^[:space:]]+)[[:space:]]+(.*)[[:space:]]*$ ]]; then
    var="${BASH_REMATCH[1]}"
    def="${BASH_REMATCH[2]}"
    eval "stamp_${var}=$(printf '%q' "$def")"
  else
    die "Malformed variable_stamps.sh output line ${line}"
  fi
done < <("${SCRIPT_DIR}/variable_stamps.sh")

while read -r line || [[ -n "$line" ]]; do
  if [[ "$line" =~ ^[[:space:]]*$ ]]; then
    continue
  elif [[ "$line" =~ ^([^:]+):([[:digit:]]+):[[:space:]]*(var[[:space:]]+)?([^[:space:]]+)[[:space:]].*//XDef:([^[:space:]]+)[[:space:]]*$ ]]; then
    go_file="${BASH_REMATCH[1]}"
    go_line="${BASH_REMATCH[2]}"
    go_var="${BASH_REMATCH[4]}"
    stamp_var="${BASH_REMATCH[5]}"

    varname="stamp_${stamp_var}"
    [[ -n "${!varname}" ]] || x_def_errors+=(
      "Variable ${go_var} defined in ${go_file}:${go_line} references status var ${stamp_var} that is not part of the variable_stamps.sh output"
    )
    go_package="$(cd "${SCRIPT_DIR}/.."; go list -e "./$(dirname "$go_file")")"

    x_defs+=(-X "\"${go_package}.${go_var}=${!varname}\"")
  fi
done < <(git -C "${SCRIPT_DIR}/.." grep -n '//XDef:' -- '*.go')
if [[ "${#x_def_errors[@]}" -gt 0 ]]; then
  printf >&2 "%s\n" "${x_def_errors[@]}"
  exit 1
fi

ldflags=(-s -w "${x_defs[@]}")

[[ -n "${GOOS}" ]] || die "GOOS must be set"
bin_name="$(basename "$main_srcdir")"
output_file="bin/${GOOS}/${bin_name}"
if [[ "$GOOS" == "windows" ]]; then
  output_file="${output_file}.exe"
fi
mkdir -p "$(dirname "$output_file")"
echo >&2 "Compiling Go source in ${main_srcdir} to ${output_file}"
go build -ldflags="${ldflags[*]}" -o "${output_file}" "${main_srcdir}"
