#!/usr/bin/env bash

set -euo pipefail

cleanup() {
  local temp_dir="$1"
  if [[ -z "${temp_dir}" ]]; then
    echo >&2 "Note: No temp dir to be deleted during cleanup."
    return
  fi
  echo "Deleting temporary directory '${temp_dir}'..."
  echo rm -rf "${temp_dir}"
}

update_crds() {
  local stackrox_ref="$1"

  # Files to copy from stackrox repo.
  CRD_FILES=(
    "central_types.go"
    "common_types.go"
    "groupversion_info.go"
    "overlay_types.go"
    "securedcluster_types.go"
    "zz_generated.deepcopy.go"
  )

  echo "Updating StackRox CRDs from ref ${stackrox_ref}..."

  # Get the directory where this script is located.
  local script_dir
  script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

  local api_version="v1alpha1"
  local target_dir="${script_dir}/${api_version}"

  # Create a temporary directory and register cleanup handler.
  local temp_dir
  temp_dir=$(mktemp -d)
  # shellcheck disable=SC2064
  trap "cleanup '${temp_dir}'" EXIT

  echo "Cloning stackrox/stackrox repository..."
  git clone https://github.com/stackrox/stackrox.git "${temp_dir}/stackrox"
  git -C "${temp_dir}/stackrox" checkout "${stackrox_ref}"

  local source_dir="${temp_dir}/stackrox/operator/api/v1alpha1"

  # Check if source directory exists.
  if [[ ! -d "${source_dir}" ]]; then
    echo "Error: Source directory ${source_dir} does not exist"
    exit 1
  fi

  for file in "${CRD_FILES[@]}"; do
    if [[ ! -f "${source_dir}/${file}" ]]; then
      echo >&2 "Error: ${file} not found in source directory -- CRD file structure must have changed."
      echo >&2 "Please update script accordingly."
      exit 1
    fi
  done

  echo "Copying CRD files..."

  for file in "${CRD_FILES[@]}"; do
    echo "Copying ${file}..."
    cp "${source_dir}/${file}" "${target_dir}/${file}"
  done

  echo "Successfully updated StackRox CRDs from tag ${stackrox_ref}"
}


# Check if STACKROX_REF set.
if [[ -z "${STACKROX_REF:-}" ]]; then
  echo >&2 "Error: STACKROX_REF environment variable must be set"
  echo >&2 "Example usage: STACKROX_REF=<git ref like tag or commit hash> ./update-crds.sh"
  exit 1
fi

update_crds "${STACKROX_REF}"
