#!/usr/bin/env bats

load bats-support-clone
load test_helper/bats-support/load
load test_helper/redhatcop-bats-library/src/error-handling

# NOTE: Each test matches to a built-in check outputted via 'kube-linter check'

get_value_from() {
  value=$(echo "${1}" | jq -r "${2}")

  if [[ -z "${value}" ]] ; then
    fail "# FATAL-ERROR: get_value_from: value is empty or invalid" || return $?
  fi

  echo "${value}"
}

@test "template-check-installed-bash-version" {
    run "bash --version"
    [[ "${BASH_VERSION:0:1}" -ge '4' ]] || false
}

@test "access-to-create-pods" {
  tmp="tests/checks/access-to-create-pods.yml"
  cmd="${KUBE_LINTER_BIN} lint --include access-to-create-pods --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "RoleBinding: binding to \"role1\" role that has [create] access to [pods]" ]]
  [[ "${count}" == "1" ]]
}

@test "access-to-secrets" {
  tmp="tests/checks/access-to-secrets.yml"
  cmd="${KUBE_LINTER_BIN} lint --include access-to-secrets --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "RoleBinding: binding to \"role1\" role that has [get] access to [secrets]" ]]
  [[ "${count}" == "1" ]]
}

@test "cluster-admin-role-binding" {
  tmp="tests/checks/cluster-admin-role-binding.yml"
  cmd="${KUBE_LINTER_BIN} lint --include cluster-admin-role-binding --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "ClusterRoleBinding: cluster-admin role is bound to [{\"kind\":\"ServiceAccount\",\"name\":\"account1\",\"namespace\":\"namespace-dev\"}]" ]]
  [[ "${count}" == "1" ]]
}

@test "dangling-horizontalpodautoscaler" {
  tmp="tests/checks/dangling-hpa.yml"
  cmd="${KUBE_LINTER_BIN} lint --include dangling-horizontalpodautoscaler --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "HorizontalPodAutoscaler: no resources found matching HorizontalPodAutoscaler scaleTargetRef ({Deployment app2 apps/v1})" ]]
  [[ "${count}" == "1" ]]
}

@test "dangling-ingress" {
  tmp="tests/checks/dangling-ingress.yml"
  cmd="${KUBE_LINTER_BIN} lint --include dangling-ingress --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Ingress: no service found matching ingress label (missing), port 80" ]]
  [[ "${message2}" == "Ingress: no service found matching ingress label (bad-port), port 8080" ]]
  [[ "${count}" == "2" ]]
}

@test "dangling-networkpolicy" {
  tmp="tests/checks/dangling-networkpolicy.yml"
  cmd="${KUBE_LINTER_BIN} lint --include dangling-networkpolicy --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "NetworkPolicy: no pods found matching networkpolicy's podSelector labels ({map[app.kubernetes.io/name:app1] []}) " ]]
  [[ "${count}" == "1" ]]
}

@test "dangling-networkpolicypeer-podselector" {
  tmp="tests/checks/dangling-networkpolicypeer-podselector.yml"
  cmd="${KUBE_LINTER_BIN} lint --include dangling-networkpolicypeer-podselector --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "NetworkPolicy: no pods found matching networkpolicy rule's podSelector labels (app.kubernetes.io/name=app2)" ]]
  [[ "${count}" == "1" ]]
}

@test "dangling-service" {
  tmp="tests/checks/dangling-service.yml"
  cmd="${KUBE_LINTER_BIN} lint --include dangling-service --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Service: no pods found matching service labels (map[app.kubernetes.io/name:app])" ]]
  [[ "${message2}" == "Service: no pods found matching service labels (map[app.kubernetes.io/name:app])" ]]
  [[ "${count}" == "2" ]]
}

@test "dangling-servicemonitor" {
  tmp="tests/checks/dangling-servicemonitor.yml"
  cmd="${KUBE_LINTER_BIN} lint --include dangling-servicemonitor --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  message3=$(get_value_from "${lines[0]}" '.Reports[2].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[2].Diagnostic.Message')
  message4=$(get_value_from "${lines[0]}" '.Reports[3].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[3].Diagnostic.Message')

  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "ServiceMonitor: no services found matching the service monitor's label selector (app.kubernetes.io/name=app) and namespace selector ([])" ]]
  [[ "${message2}" == "ServiceMonitor: no services found matching the service monitor's label selector (app.kubernetes.io/name=app) and namespace selector ([])" ]]
  [[ "${message3}" == "ServiceMonitor: no services found matching the service monitor's label selector () and namespace selector ([test2])" ]]
  [[ "${message4}" == "ServiceMonitor: no services found matching the service monitor's label selector (app.kubernetes.io/name=app1) and namespace selector ([test2])" ]]
  [[ "${count}" == "4" ]]
}

@test "default-service-account" {
  tmp="tests/checks/default-service-account.yml"
  cmd="${KUBE_LINTER_BIN} lint --include default-service-account --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: found matching serviceAccount (\"default\")" ]]
  [[ "${message2}" == "DeploymentConfig: found matching serviceAccount (\"default\")" ]]
  [[ "${count}" == "2" ]]
}

@test "deprecated-service-account-field" {
  tmp="tests/checks/deprecated-service-account-field.yml"
  cmd="${KUBE_LINTER_BIN} lint --include deprecated-service-account-field --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: serviceAccount is specified (default), but this field is deprecated; use serviceAccountName instead" ]]
  [[ "${message2}" == "DeploymentConfig: serviceAccount is specified (default), but this field is deprecated; use serviceAccountName instead" ]]
  [[ "${count}" == "2" ]]
}

@test "dnsconfig-options" {
  tmp="tests/checks/dnsconfig-options-ndots.yml"
  cmd="${KUBE_LINTER_BIN} lint --include dnsconfig-options --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  message3=$(get_value_from "${lines[0]}" '.Reports[2].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[2].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: DNSConfig Options \"ndots:2\" not found." ]]
  [[ "${message2}" == "Deployment: Object does not define any DNSConfig Options." ]]
  [[ "${message3}" == "Deployment: Object does not define any DNSConfig rules." ]]
  [[ "${count}" == "3" ]]
}

@test "docker-sock" {
  tmp="tests/checks/docker-sock.yml"
  cmd="${KUBE_LINTER_BIN} lint --include docker-sock --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: host system directory \"/var/run/docker.sock\" is mounted on container \"app\"" ]]
  [[ "${message2}" == "DeploymentConfig: host system directory \"/var/run/docker.sock\" is mounted on container \"app\"" ]]
  [[ "${count}" == "2" ]]
}

@test "drop-net-raw-capability" {
  tmp="tests/checks/drop-net-raw-capability.yml"
  cmd="${KUBE_LINTER_BIN} lint --include drop-net-raw-capability --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[2].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[2].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: container \"app\" has ADD capability: \"NET_RAW\", which matched with the forbidden capability for containers" ]]
  [[ "${message2}" == "DeploymentConfig: container \"app\" has ADD capability: \"NET_RAW\", which matched with the forbidden capability for containers" ]]
  [[ "${count}" == "4" ]]
}

@test "duplicate-env-var" {
  tmp="tests/checks/duplicate-env-var.yaml"
  cmd="${KUBE_LINTER_BIN} lint --include duplicate-env-var --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: Duplicate environment variable PORT in container \"fire-deployment\" found" ]]
  [[ "${message2}" == "StatefulSet: Duplicate environment variable PORT in container \"fire-stateful\" found" ]]
  [[ "${count}" == "2" ]]
}

@test "env-var-secret" {
  tmp="tests/checks/env-var-secret.yml"
  cmd="${KUBE_LINTER_BIN} lint --include env-var-secret --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: environment variable SECRET_BLAH in container \"app\" found" ]]
  [[ "${message2}" == "DeploymentConfig: environment variable SECRET_BLAH in container \"app\" found" ]]
  [[ "${count}" == "2" ]]
}

@test "exposed-services" {
  tmp="tests/checks/exposed-services.yml"
  cmd="${KUBE_LINTER_BIN} lint --include exposed-services --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Service: \"NodePort\" service type is forbidden." ]]
  [[ "${count}" == "1" ]]
}

@test "host-ipc" {
  tmp="tests/checks/host-ipc.yml"
  cmd="${KUBE_LINTER_BIN} lint --include host-ipc --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: resource shares host's IPC namespace (via hostIPC=true)." ]]
  [[ "${message2}" == "DeploymentConfig: resource shares host's IPC namespace (via hostIPC=true)." ]]
  [[ "${count}" == "2" ]]
}

@test "host-network" {
  tmp="tests/checks/host-network.yml"
  cmd="${KUBE_LINTER_BIN} lint --include host-network --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: resource shares host's network namespace (via hostNetwork=true)." ]]
  [[ "${message2}" == "DeploymentConfig: resource shares host's network namespace (via hostNetwork=true)." ]]
  [[ "${count}" == "2" ]]
}

@test "host-pid" {
  tmp="tests/checks/host-pid.yml"
  cmd="${KUBE_LINTER_BIN} lint --include host-pid --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: object shares the host's process namespace (via hostPID=true)." ]]
  [[ "${message2}" == "DeploymentConfig: object shares the host's process namespace (via hostPID=true)." ]]
  [[ "${count}" == "2" ]]
}

@test "hpa-minimum-three-replicas" {
  tmp="tests/checks/hpa-minimum-three-replicas.yml"
  cmd="${KUBE_LINTER_BIN} lint --include hpa-minimum-three-replicas --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "HorizontalPodAutoscaler: object has 2 replicas but minimum required replicas is 3" ]]
  [[ "${count}" == "1" ]]
}

@test "invalid-target-ports" {
  tmp="tests/checks/invalid-target-ports.yaml"
  cmd="${KUBE_LINTER_BIN} lint --include invalid-target-ports --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [[ "${status}" -eq 1 ]]

  count=$(get_value_from "${lines[0]}" '.Reports | length')
  [[ "${count}" == "4" ]]

  # TODO: export to helper function so that it can be used for other tests
  actual_messages=()
  for (( i=0; i<$((count)); i++ ))
  do
    actual_message=$(get_value_from "${lines[0]}" ".Reports[${i}] | .Object.K8sObject.GroupVersionKind.Kind + \": \" + .Diagnostic.Message")
    actual_messages+=("${actual_message}")
  done

  [[ "${actual_messages[0]}" == "Service: port targetPort \"123456\" in service \"invalid-target-ports\" must be between 1 and 65535, inclusive" ]]
  [[ "${actual_messages[1]}" == "Service: port targetPort \"n234567890123456\" in service \"invalid-target-ports\" must be no more than 15 characters" ]]
  [[ "${actual_messages[2]}" == "Deployment: port name \"n234567890123456\" in container \"invalid-target-ports\" must be no more than 15 characters" ]]
  [[ "${actual_messages[3]}" == "Deployment: port name \"123456\" in container \"invalid-target-ports\" must contain at least one letter (a-z)" ]]
}

@test "latest-tag" {
  tmp="tests/checks/latest-tag.yml"
  cmd="${KUBE_LINTER_BIN} lint --include latest-tag --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: The container \"app\" is using an invalid container image, \"app:latest\". Please use images that are not blocked by the \`BlockList\` criteria : [\".*:(latest)$\" \"^[^:]*$\" \"(.*/[^:]+)$\"]" ]]
  [[ "${message2}" == "DeploymentConfig: The container \"app\" is using an invalid container image, \"app:latest\". Please use images that are not blocked by the \`BlockList\` criteria : [\".*:(latest)$\" \"^[^:]*$\" \"(.*/[^:]+)$\"]" ]]
  [[ "${count}" == "2" ]]
}

@test "liveness-port" {
  tmp="tests/checks/liveness-port.yml"
  cmd="${KUBE_LINTER_BIN} lint --include liveness-port --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  message3=$(get_value_from "${lines[0]}" '.Reports[2].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[2].Diagnostic.Message')
  message4=$(get_value_from "${lines[0]}" '.Reports[3].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[3].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: container \"fire-deployment-name\" does not expose port http for the HTTPGet" ]]
  [[ "${message2}" == "Deployment: container \"fire-deployment-int\" does not expose port 8080 for the HTTPGet" ]]
  [[ "${message3}" == "Deployment: container \"fire-deployment-udp\" does not expose port udp for the TCPSocket" ]]
  [[ "${message4}" == "StatefulSet: container \"fire-stateful-name\" does not expose port healthcheck for the HTTPGet" ]]
  [[ "${count}" == "4" ]]
}


@test "minimum-three-replicas" {
  tmp="tests/checks/minimum-three-replicas.yml"
  cmd="${KUBE_LINTER_BIN} lint --include minimum-three-replicas --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: object has 2 replicas but minimum required replicas is 3" ]]
  [[ "${message2}" == "DeploymentConfig: object has 2 replicas but minimum required replicas is 3" ]]
  [[ "${count}" == "2" ]]
}

@test "mismatching-selector" {
  tmp="tests/checks/mismatching-selector.yml"
  cmd="${KUBE_LINTER_BIN} lint --include mismatching-selector --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: labels in pod spec (map[]) do not match labels in selector (&LabelSelector{MatchLabels:map[string]string{app.kubernetes.io/name: app1,},MatchExpressions:[]LabelSelectorRequirement{},})" ]]
  [[ "${message2}" == "DeploymentConfig: labels in pod spec (map[]) do not match labels in selector (&LabelSelector{MatchLabels:map[string]string{app.kubernetes.io/name: app2,},MatchExpressions:[]LabelSelectorRequirement{},})" ]]
  [[ "${count}" == "2" ]]
}

@test "no-anti-affinity" {
  tmp="tests/checks/no-anti-affinity.yml"
  cmd="${KUBE_LINTER_BIN} lint --include no-anti-affinity --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " +.Reports[1].Diagnostic.Message')
  message3=$(get_value_from "${lines[0]}" '.Reports[2].Object.K8sObject.GroupVersionKind.Kind + ": " +.Reports[2].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: object has 3 replicas but does not specify inter pod anti-affinity" ]]
  [[ "${message2}" == "DeploymentConfig: object has 3 replicas but does not specify inter pod anti-affinity" ]]
  [[ "${message3}" == "Deployment: pod's namespace \"foo\" not found in anti-affinity's namespaces [bar]" ]]
  [[ "${count}" == "3" ]]
}

@test "no-extensions-v1beta" {
  tmp="tests/checks/no-extensions-v1beta.yml"
  cmd="${KUBE_LINTER_BIN} lint --include no-extensions-v1beta --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "NetworkPolicy: disallowed API object found: extensions/v1beta1, Kind=NetworkPolicy" ]]
  [[ "${count}" == "1" ]]
}

@test "no-liveness-probe" {
  tmp="tests/checks/no-liveness-probe.yml"
  cmd="${KUBE_LINTER_BIN} lint --include no-liveness-probe --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: container \"app1\" does not specify a liveness probe" ]]
  [[ "${message2}" == "DeploymentConfig: container \"app2\" does not specify a liveness probe" ]]
  [[ "${count}" == "2" ]]
}

@test "no-node-affinity" {
  tmp="tests/checks/no-node-affinity.yml"
  cmd="${KUBE_LINTER_BIN} lint --include no-node-affinity --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: object does not define any node affinity rules." ]]
  [[ "${count}" == "1" ]]
}

@test "no-read-only-root-fs" {
  tmp="tests/checks/no-read-only-root-fs.yml"
  cmd="${KUBE_LINTER_BIN} lint --include no-read-only-root-fs --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: container \"app\" does not have a read-only root file system" ]]
  [[ "${message2}" == "DeploymentConfig: container \"app\" does not have a read-only root file system" ]]
  [[ "${count}" == "2" ]]
}

@test "no-readiness-probe" {
  tmp="tests/checks/no-readiness-probe.yml"
  cmd="${KUBE_LINTER_BIN} lint --include no-readiness-probe --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: container \"app\" does not specify a readiness probe" ]]
  [[ "${message2}" == "DeploymentConfig: container \"app\" does not specify a readiness probe" ]]
  [[ "${count}" == "2" ]]
}

@test "no-rolling-update-strategy" {
  tmp="tests/checks/no-rolling-update-strategy.yml"
  cmd="${KUBE_LINTER_BIN} lint --include no-rolling-update-strategy --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: object has Other strategy type but must match regex ^(RollingUpdate|Rolling)$" ]]
  [[ "${message2}" == "DeploymentConfig: object has Other strategy type but must match regex ^(RollingUpdate|Rolling)$" ]]
  [[ "${count}" == "2" ]]
}

@test "non-existent-service-account" {
  tmp="tests/checks/non-existent-service-account.yml"
  cmd="${KUBE_LINTER_BIN} lint --include non-existent-service-account --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: serviceAccount \"missing\" not found" ]]
  [[ "${message2}" == "DeploymentConfig: serviceAccount \"missing\" not found" ]]
  [[ "${count}" == "2" ]]
}

@test "non-isolated-pod" {
  tmp="tests/checks/non-isolated-pod.yml"
  cmd="${KUBE_LINTER_BIN} lint --include non-isolated-pod --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: pods created by this object are non-isolated" ]]
  [[ "${message2}" == "DeploymentConfig: pods created by this object are non-isolated" ]]
  [[ "${count}" == "2" ]]
}

@test "pdb-max-unavailable" {

  tmp="tests/checks/pdb-max-unavailable.yaml"
  cmd="${KUBE_LINTER_BIN} lint --include pdb-max-unavailable --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')

  [[ "${message1}" == "PodDisruptionBudget: MaxUnavailable is set to 0" ]]

}

@test "pdb-min-available" {

  tmp="tests/checks/pdb-min-available.yaml"
  cmd="${KUBE_LINTER_BIN} lint --include pdb-min-available --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}


  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  message3=$(get_value_from "${lines[0]}" '.Reports[2].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[2].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "PodDisruptionBudget: The current number of replicas for deployment foo is equal to or lower than the minimum number of replicas specified by its PDB." ]]
  [[ "${message2}" == "PodDisruptionBudget: The current number of replicas for deployment foo2 is equal to or lower than the minimum number of replicas specified by its PDB." ]]
  [[ "${message3}" == "PodDisruptionBudget: The current number of replicas for deployment foo3 is equal to or lower than the minimum number of replicas specified by its PDB." ]]
  [[ "${count}" == "3" ]]
}

@test "pdb-unhealthy-pod-eviction-policy" {

  tmp="tests/checks/pdb-unhealthy-pod-eviction-policy.yaml"
  cmd="${KUBE_LINTER_BIN} lint --include pdb-unhealthy-pod-eviction-policy --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')

  [[ "${message1}" == "PodDisruptionBudget: unhealthyPodEvictionPolicy is not explicitly set" ]]
  count=$(get_value_from "${lines[0]}" '.Reports | length')
  [[ "${count}" == "1" ]]

}

@test "privilege-escalation-container" {
  tmp="tests/checks/privilege-escalation-container.yml"
  cmd="${KUBE_LINTER_BIN} lint --include privilege-escalation-container --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: container \"app\" has AllowPrivilegeEscalation set to true." ]]
  [[ "${message2}" == "DeploymentConfig: container \"app\" has AllowPrivilegeEscalation set to true." ]]
  [[ "${count}" == "2" ]]
}

@test "privileged-container" {
  tmp="tests/checks/privileged-container.yml"
  cmd="${KUBE_LINTER_BIN} lint --include privileged-container --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: container \"app\" is privileged" ]]
  [[ "${message2}" == "DeploymentConfig: container \"app\" is privileged" ]]
  [[ "${count}" == "2" ]]
}

@test "privileged-ports" {
  tmp="tests/checks/privileged-ports.yml"
  cmd="${KUBE_LINTER_BIN} lint --include privileged-ports --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: port 80 is mapped in container \"app\"." ]]
  [[ "${message2}" == "DeploymentConfig: port 80 is mapped in container \"app\"." ]]
  [[ "${count}" == "2" ]]
}

@test "read-secret-from-env-var" {
  tmp="tests/checks/read-secret-from-env-var.yml"
  cmd="${KUBE_LINTER_BIN} lint --include read-secret-from-env-var --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: environment variable \"TOKEN\" in container \"app\" uses SecretKeyRef" ]]
  [[ "${message2}" == "DeploymentConfig: environment variable \"TOKEN\" in container \"app\" uses SecretKeyRef" ]]
  [[ "${count}" == "2" ]]
}

@test "readiness-port" {
  tmp="tests/checks/readiness-port.yml"
  cmd="${KUBE_LINTER_BIN} lint --include readiness-port --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  message3=$(get_value_from "${lines[0]}" '.Reports[2].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[2].Diagnostic.Message')
  message4=$(get_value_from "${lines[0]}" '.Reports[3].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[3].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: container \"fire-deployment-name\" does not expose port http for the HTTPGet" ]]
  [[ "${message2}" == "Deployment: container \"fire-deployment-int\" does not expose port 8080 for the HTTPGet" ]]
  [[ "${message3}" == "Deployment: container \"fire-deployment-udp\" does not expose port udp for the TCPSocket" ]]
  [[ "${message4}" == "Deployment: container \"fire-deployment-grpc\" does not expose port 8080 for the GRPC check" ]]
  [[ "${count}" == "4" ]]
}

@test "required-annotation-email" {
  tmp="tests/checks/required-annotation-email.yml"
  cmd="${KUBE_LINTER_BIN} lint --include required-annotation-email --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: no annotation matching \"email=[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+\" found" ]]
  [[ "${message2}" == "DeploymentConfig: no annotation matching \"email=[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+\" found" ]]
  [[ "${count}" == "2" ]]
}

@test "required-label-owner" {
  tmp="tests/checks/required-label-owner.yml"
  cmd="${KUBE_LINTER_BIN} lint --include required-label-owner --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: no label matching \"owner=<any>\" found" ]]
  [[ "${message2}" == "DeploymentConfig: no label matching \"owner=<any>\" found" ]]
  [[ "${count}" == "2" ]]
}

@test "run-as-non-root" {
  tmp="tests/checks/run-as-non-root.yml"
  cmd="${KUBE_LINTER_BIN} lint --include run-as-non-root --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: container \"app\" is not set to runAsNonRoot" ]]
  [[ "${message2}" == "DeploymentConfig: container \"app2\" is not set to runAsNonRoot" ]]
  [[ "${count}" == "2" ]]
}

@test "scc-deny-privileged-container" {
  tmp="tests/checks/scc-deny-privileged-container.yml"
  cmd="${KUBE_LINTER_BIN} lint --include scc-deny-privileged-container --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "SecurityContextConstraints: SCC has allowPrivilegedContainer set to true" ]]
  [[ "${count}" == "1" ]]
}

@test "sensitive-host-mounts" {
  tmp="tests/checks/sensitive-host-mounts.yml"
  cmd="${KUBE_LINTER_BIN} lint --include sensitive-host-mounts --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: host system directory \"/etc\" is mounted on container \"app\"" ]]
  [[ "${message2}" == "DeploymentConfig: host system directory \"/etc\" is mounted on container \"app2\"" ]]
  [[ "${count}" == "2" ]]
}

@test "ssh-port" {
  tmp="tests/checks/ssh-port.yml"
  cmd="${KUBE_LINTER_BIN} lint --include ssh-port --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  message3=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[2].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: port 22 and protocol TCP in container \"app\" found" ]]
  [[ "${message2}" == "DeploymentConfig: port 22 and protocol TCP in container \"app\" found" ]]
  [[ "${message3}" == "DeploymentConfig: port 22 and protocol TCP in container \"app-no-protocol\" found" ]]
  [[ "${count}" == "3" ]]
}

@test "startup-port" {
  tmp="tests/checks/startup-port.yml"
  cmd="${KUBE_LINTER_BIN} lint --include startup-port --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  message3=$(get_value_from "${lines[0]}" '.Reports[2].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[2].Diagnostic.Message')
  message4=$(get_value_from "${lines[0]}" '.Reports[3].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[3].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: container \"fire-deployment-name\" does not expose port http for the HTTPGet" ]]
  [[ "${message2}" == "Deployment: container \"fire-deployment-int\" does not expose port 8080 for the HTTPGet" ]]
  [[ "${message3}" == "Deployment: container \"fire-deployment-udp\" does not expose port udp for the TCPSocket" ]]
  [[ "${message4}" == "Deployment: container \"fire-deployment-grpc\" does not expose port 8080 for the GRPC check" ]]
  [[ "${count}" == "4" ]]
}

@test "unsafe-proc-mount" {
  tmp="tests/checks/unsafe-proc-mount.yml"
  cmd="${KUBE_LINTER_BIN} lint --include unsafe-proc-mount --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: container \"app\" exposes /proc unsafely (via procMount=Unmasked)." ]]
  [[ "${message2}" == "DeploymentConfig: container \"app\" exposes /proc unsafely (via procMount=Unmasked)." ]]
  [[ "${count}" == "2" ]]
}

@test "unsafe-sysctls" {
  tmp="tests/checks/unsafe-sysctls.yml"
  cmd="${KUBE_LINTER_BIN} lint --include unsafe-sysctls --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: resource specifies unsafe sysctl \"kernel.sem\"." ]]
  [[ "${message2}" == "DeploymentConfig: resource specifies unsafe sysctl \"kernel.sem\"." ]]
  [[ "${count}" == "2" ]]
}

@test "unset-cpu-requirements" {
  tmp="tests/checks/unset-cpu-requirements.yml"
  cmd="${KUBE_LINTER_BIN} lint --include unset-cpu-requirements --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: container \"app\" has cpu request 0" ]]
  [[ "${message2}" == "DeploymentConfig: container \"app\" has cpu request 0" ]]
  [[ "${count}" == "2" ]]
}

@test "unset-memory-requirements" {
  tmp="tests/checks/unset-memory-requirements.yml"
  cmd="${KUBE_LINTER_BIN} lint --include unset-memory-requirements --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: container \"app\" has memory limit 0" ]]
  [[ "${message2}" == "DeploymentConfig: container \"app\" has memory limit 0" ]]
  [[ "${count}" == "2" ]]
}

@test "use-namespace" {
  tmp="tests/checks/use-namespace.yml"
  cmd="${KUBE_LINTER_BIN} lint --include use-namespace --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: object in default namespace" ]]
  [[ "${count}" == "1" ]]
}

@test "wildcard-in-rules" {
  tmp="tests/checks/wildcard-in-rules.yml"
  cmd="${KUBE_LINTER_BIN} lint --include wildcard-in-rules --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Role: wildcard "*" in verb specification" ]]
  [[ "${count}" == "1" ]]
}

@test "writable-host-mount" {
  tmp="tests/checks/writable-host-mount.yml"
  cmd="${KUBE_LINTER_BIN} lint --include writable-host-mount --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  message2=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[1].Diagnostic.Message')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: container app mounts path /config on the host as writable" ]]
  [[ "${message2}" == "DeploymentConfig: container app2 mounts path /config on the host as writable" ]]
  [[ "${count}" == "2" ]]
}

@test "template-forbidden-annotation" {
  tmp="tests/checks/forbidden-annotation.yml"
  cmd="${KUBE_LINTER_BIN} lint --config e2etests/testdata/forbidden-annotation-config.yaml --do-not-auto-add-defaults --format json ${tmp}"
  run ${cmd}

  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 1 ]

  message1=$(get_value_from "${lines[0]}" '.Reports[0].Object.K8sObject.GroupVersionKind.Kind + ": " + .Reports[0].Diagnostic.Message')
  failing_resource=$(get_value_from "${lines[0]}" '.Reports[1].Object.K8sObject.Name')
  count=$(get_value_from "${lines[0]}" '.Reports | length')

  [[ "${message1}" == "Deployment: annotation matching \"reloader.stakater.com/auto=true\" found" ]]
  [[ "${failing_resource}" == "bad-irsa-role" ]]
  [[ "${count}" == "2" ]]
}

@test "flag-ignore-paths" {
  tmp="."
  cmd="${KUBE_LINTER_BIN} lint --ignore-paths \"tests/**\" --ignore-paths \"e2etests/**\" ${tmp}"
  run ${cmd}
  print_info "${status}" "${output}" "${cmd}" "${tmp}"
  [ "$status" -eq 0 ]
}

@test "flag-read-from-stdin" {
  echo "---" | ${KUBE_LINTER_BIN} lint -
}
