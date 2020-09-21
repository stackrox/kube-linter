#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "${SCRIPT_DIR}/utils.sh"

set -euo pipefail

gitroot="$(git rev-parse --show-toplevel)"

[[ -n "${gitroot}" ]] || die "Could not determine git root"

echo "VERSION $("${gitroot}/get-tag")"
