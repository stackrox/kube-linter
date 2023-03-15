#!/usr/bin/env bash

run() {
  local tmp_write_dir=/tmp/kubelinter/$(date +'%d-%m-%Y-%H-%M')
  mkdir -p "${tmp_write_dir}"

  grep "@test" e2etests/bats-tests.sh | grep -v 'flag-' | grep -v 'template-' | cut -d'"' -f2 > ${tmp_write_dir}/batstests.log
  ${KUBE_LINTER_BIN:-kube-linter} checks list --format json | jq -r '.[].name' > ${tmp_write_dir}/kubelinterchecks.log
  diff -c ${tmp_write_dir}/kubelinterchecks.log ${tmp_write_dir}/batstests.log || { echo >&2 "ERROR: The output of '${KUBE_LINTER_BIN} checks list' differs from the tests in 'e2etests/bats-tests.sh'. See above diff."; exit 1; }
}

run
